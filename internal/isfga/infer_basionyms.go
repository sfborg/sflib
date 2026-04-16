package isfga

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sfborg/sflib/pkg/coldp"
)

// nameForInference holds the data needed for basionym matching.
type nameForInference struct {
	ID                    string
	CanonicalStemmed      string
	Authorship            string // full verbatim authorship (may contain parentheses)
	OriginalAuthorship    string // author without parentheses (basionym author)
	CombinationAuthorship string // author in parentheses
	Year                  string // basionym year or published_in_year
	Genus                 string
	InfragenericEpithet   string // subgenus
	SpecificEpithet       string
	InfraspecificEpithet  string
	Rank                  string
}

// nameLookups holds pre-built maps for fast name ID lookups.
// This avoids N+1 queries when creating OriginalX relationships.
type nameLookups struct {
	// genus: uninomial -> ID (for genus-rank names)
	genus map[string]string
	// subgenus: "genus|infrageneric" -> ID
	subgenus map[string]string
	// species: "genus|specific" -> ID
	species map[string]string
	// infraspecific: "genus|specific|infraspecific|rank" -> ID
	infraspecific map[string]string
}

// InferBasionyms detects and creates basionym relationships by matching
// stemmed epithets and original authorship across all names in the archive.
//
// Algorithm:
// 1. Build lookup map: (stemmed_epithet + original_authorship + year) -> name_id
// 2. Names without parentheses in authorship are potential basionyms
// 3. Names with parenthetical authorship are combinations (genus transfer)
// 4. Match combinations to basionyms; skip ambiguous matches (multiple candidates)
// 5. Create BASIONYM name_relation records for matches
// 6. Optionally create OriginalX relationships (OriginalGenus, OriginalSubgenus,
//    OriginalSpecies, OriginalSubspecies, OriginalVariety, etc.)
//
// Key matching uses: stemmed_lowest_epithet + basionym_authorship + year
// Year helps disambiguate when same author described same epithet in different years.
// Parentheses detection checks both col__combination_authorship AND parses col__authorship.
func (a *isfga) InferBasionyms(ctx context.Context) error {
	// Check if we should skip due to existing relations
	if a.cfg.SkipBasionymsIfRelationsExist {
		exists, err := a.basionymRelationsExist()
		if err != nil {
			return fmt.Errorf("checking existing relations: %w", err)
		}
		if exists {
			slog.Info("Skipping basionym inference: BASIONYM relations already exist")
			return nil
		}
	}

	slog.Info("Starting basionym inference")

	// Load all names with their parsed data
	names, err := a.loadNamesForInference(ctx)
	if err != nil {
		return fmt.Errorf("loading names: %w", err)
	}
	slog.Info("Loaded names for inference", "count", len(names))

	// Build the basionym lookup and blacklist
	basionymLookup := make(map[string]nameForInference)
	blacklist := make(map[string]struct{})

	for _, n := range names {
		// Only names without combination authorship can be basionyms
		if !isBasionym(n) {
			continue
		}

		key := getBasionymKey(n)
		if key == "" {
			continue
		}

		if _, exists := blacklist[key]; exists {
			// Already blacklisted as ambiguous
			continue
		}

		if existing, exists := basionymLookup[key]; exists {
			// Ambiguous: multiple potential basionyms with same key
			slog.Debug("Ambiguous basionym key, blacklisting",
				"key", key,
				"name1", existing.ID,
				"name2", n.ID,
			)
			blacklist[key] = struct{}{}
			delete(basionymLookup, key)
		} else {
			basionymLookup[key] = n
		}
	}

	slog.Info("Built basionym lookup",
		"candidates", len(basionymLookup),
		"blacklisted", len(blacklist),
	)

	// Build lookup maps for OriginalX relationships (avoids N+1 queries)
	var lookups nameLookups
	if a.cfg.CreateOriginalCombinations {
		slog.Info("Building name lookup maps for OriginalX relationships")
		var err error
		lookups, err = a.buildNameLookups(ctx)
		if err != nil {
			return fmt.Errorf("building name lookups: %w", err)
		}
		slog.Info("Built name lookups",
			"genera", len(lookups.genus),
			"subgenera", len(lookups.subgenus),
			"species", len(lookups.species),
			"infraspecific", len(lookups.infraspecific),
		)
	}

	// Match combinations to basionyms
	var relations []coldp.NameRelation
	matchCount := 0

	for _, n := range names {
		// Only names with combination authorship (parentheses) are combinations
		if !isCombination(n) {
			continue
		}
		if n.OriginalAuthorship == "" {
			continue
		}

		key := getBasionymKey(n)
		if key == "" {
			continue
		}

		basionym, exists := basionymLookup[key]
		if !exists {
			continue
		}

		// Don't match to self
		if basionym.ID == n.ID {
			continue
		}

		matchCount++

		// Create BASIONYM relation: combination -> basionym
		// The tw_type reflects the BASIONYM's rank, not the combination's rank.
		// This tells us what level the original name was at.
		twType := inferOriginalTypeFromRank(basionym.Rank)
		relations = append(relations, coldp.NameRelation{
			NameID:                 n.ID,
			RelatedNameID:          basionym.ID,
			Type:                   coldp.NewNomRelType("BASIONYM"),
			TwNameRelationshipType: twType,
		})

		// Optionally create OriginalCombination relationships
		if a.cfg.CreateOriginalCombinations {
			origRels := createOriginalCombinationRelations(n, basionym, lookups)
			relations = append(relations, origRels...)
		}
	}

	slog.Info("Matched combinations to basionyms", "matches", matchCount)

	// Also create self-referencing relationships for basionyms (protonyms)
	// In both COLDP and TaxonWorks, a basionym declares itself as its own basionym
	if a.cfg.CreateOriginalCombinations {
		basionymCount := 0
		for _, basionym := range basionymLookup {
			basionymCount++

			// Self-referencing BASIONYM relation
			twType := inferOriginalTypeFromRank(basionym.Rank)
			relations = append(relations, coldp.NameRelation{
				NameID:                 basionym.ID,
				RelatedNameID:          basionym.ID,
				Type:                   coldp.NewNomRelType("BASIONYM"),
				TwNameRelationshipType: twType,
			})

			// OriginalX relationships for the basionym itself
			origRels := createOriginalCombinationRelations(basionym, basionym, lookups)
			relations = append(relations, origRels...)
		}
		slog.Info("Created self-referencing relations for basionyms", "count", basionymCount)
	}

	if len(relations) == 0 {
		slog.Info("No new basionym relations to create")
		return nil
	}

	// Insert the relations
	if err := a.InsertNameRelations(relations); err != nil {
		return fmt.Errorf("inserting basionym relations: %w", err)
	}

	slog.Info("Created basionym relations", "count", len(relations))
	return nil
}

// getBasionymKey creates a lookup key from stemmed epithet + original authorship + year.
// Including year helps disambiguate when the same author described multiple taxa
// with the same epithet in different years.
func getBasionymKey(n nameForInference) string {
	if n.CanonicalStemmed == "" || n.OriginalAuthorship == "" {
		return ""
	}

	// Get the lowest epithet from the stemmed canonical
	parts := strings.Split(n.CanonicalStemmed, " ")
	if len(parts) == 0 {
		return ""
	}
	lowestEpithet := parts[len(parts)-1]

	// Normalize authorship for comparison
	authorship := strings.TrimSpace(n.OriginalAuthorship)

	// Include year if available for better disambiguation
	if n.Year != "" {
		return lowestEpithet + "_" + authorship + "_" + n.Year
	}
	return lowestEpithet + "_" + authorship
}

// isCombination checks if a name is a combination (moved to different genus).
// Combinations have parenthetical authorship indicating the original author.
// We check both the explicit combination_authorship field AND parse
// parentheses from the full authorship string.
func isCombination(n nameForInference) bool {
	// First check the explicit field
	if n.CombinationAuthorship != "" {
		return true
	}
	// Fall back to checking for parentheses in full authorship
	// e.g., "(Smith, 1900)" or "(Smith, 1900) Jones, 1950"
	return strings.Contains(n.Authorship, "(")
}

// isBasionym checks if a name could be a basionym (original description).
// Basionyms have authorship without parentheses.
func isBasionym(n nameForInference) bool {
	if n.OriginalAuthorship == "" {
		return false
	}
	return !isCombination(n)
}

// basionymRelationsExist checks if any BASIONYM type relations exist.
func (a *isfga) basionymRelationsExist() (bool, error) {
	var count int
	err := a.db.QueryRow(`
		SELECT COUNT(*) FROM name_relation WHERE col__type_id = 'BASIONYM'
	`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// loadNamesForInference loads all names with the fields needed for basionym matching.
func (a *isfga) loadNamesForInference(ctx context.Context) ([]nameForInference, error) {
	rows, err := a.db.QueryContext(ctx, `
		SELECT
			col__id,
			COALESCE(gn__canonical_stemmed, ''),
			COALESCE(col__authorship, ''),
			COALESCE(col__basionym_authorship, ''),
			COALESCE(col__combination_authorship, ''),
			COALESCE(col__basionym_authorship_year, COALESCE(col__published_in_year, '')),
			COALESCE(col__genus, ''),
			COALESCE(col__infrageneric_epithet, ''),
			COALESCE(col__specific_epithet, ''),
			COALESCE(col__infraspecific_epithet, ''),
			COALESCE(col__rank_id, '')
		FROM name
		WHERE gn__canonical_stemmed IS NOT NULL
		  AND gn__canonical_stemmed != ''
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []nameForInference
	for rows.Next() {
		var n nameForInference
		err := rows.Scan(
			&n.ID,
			&n.CanonicalStemmed,
			&n.Authorship,
			&n.OriginalAuthorship,
			&n.CombinationAuthorship,
			&n.Year,
			&n.Genus,
			&n.InfragenericEpithet,
			&n.SpecificEpithet,
			&n.InfraspecificEpithet,
			&n.Rank,
		)
		if err != nil {
			return nil, err
		}
		names = append(names, n)
	}

	return names, rows.Err()
}

// buildNameLookups loads all names and builds lookup maps for fast ID resolution.
// This avoids N+1 queries when creating OriginalX relationships.
func (a *isfga) buildNameLookups(ctx context.Context) (nameLookups, error) {
	lookups := nameLookups{
		genus:         make(map[string]string),
		subgenus:      make(map[string]string),
		species:       make(map[string]string),
		infraspecific: make(map[string]string),
	}

	rows, err := a.db.QueryContext(ctx, `
		SELECT
			col__id,
			COALESCE(col__uninomial, ''),
			COALESCE(col__genus, ''),
			COALESCE(col__infrageneric_epithet, ''),
			COALESCE(col__specific_epithet, ''),
			COALESCE(col__infraspecific_epithet, ''),
			LOWER(COALESCE(col__rank_id, ''))
		FROM name
	`)
	if err != nil {
		return lookups, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, uninomial, genus, infrageneric, specific, infraspecific, rank string
		err := rows.Scan(&id, &uninomial, &genus, &infrageneric, &specific, &infraspecific, &rank)
		if err != nil {
			return lookups, err
		}

		switch rank {
		case "genus":
			if uninomial != "" {
				lookups.genus[uninomial] = id
			}
		case "subgenus":
			if genus != "" && infrageneric != "" {
				key := genus + "|" + infrageneric
				lookups.subgenus[key] = id
			}
		case "species":
			if genus != "" && specific != "" {
				key := genus + "|" + specific
				lookups.species[key] = id
			}
		case "subspecies", "variety", "subvariety", "form", "subform":
			if genus != "" && specific != "" && infraspecific != "" {
				key := genus + "|" + specific + "|" + infraspecific + "|" + rank
				lookups.infraspecific[key] = id
			}
		}
	}

	return lookups, rows.Err()
}

// createOriginalCombinationRelations creates OriginalGenus, OriginalSubgenus,
// OriginalSpecies, OriginalSubspecies, OriginalVariety, etc. relationships
// based on the original structure of the basionym.
//
// IMPORTANT: All OriginalX relationships point to the actual name records at
// each rank, NOT to the basionym. Each relationship captures what the original
// name was at that rank level.
//
// Handles rank changes in BOTH directions within species-group:
// - Elevation: subspecies/variety/form → species
// - Demotion: species → subspecies/variety/form
// - Lateral: subspecies → variety, etc.
//
// Example 1 (elevation): *Aus cus dus* Smith, 1900 (subspecies)
// → *Mus dus* (Smith, 1900) Jones, 1950 (species)
// Creates:
//   - OriginalGenus → "Aus"
//   - OriginalSpecies → "Aus cus"
//   - OriginalSubspecies → "Aus cus dus"
//
// Example 2 (demotion): *Aus bus* Smith, 1900 (species)
// → *Xus cus bus* (Smith, 1900) Jones, 1950 (subspecies)
// Creates:
//   - OriginalGenus → "Aus"
//   - OriginalSpecies → "Aus bus"
func createOriginalCombinationRelations(
	combination, basionym nameForInference,
	lookups nameLookups,
) []coldp.NameRelation {
	var relations []coldp.NameRelation

	// OriginalGenus: always create to document the original genus
	if basionym.Genus != "" {
		if genusID, ok := lookups.genus[basionym.Genus]; ok {
			relations = append(relations, coldp.NameRelation{
				NameID:                 combination.ID,
				RelatedNameID:          genusID,
				Type:                   coldp.NewNomRelType(""), // No COLDP equivalent
				TwNameRelationshipType: "TaxonNameRelationship::OriginalCombination::OriginalGenus",
			})
		}
	}

	// OriginalSubgenus: create if basionym had a subgenus
	if basionym.InfragenericEpithet != "" {
		key := basionym.Genus + "|" + basionym.InfragenericEpithet
		if subgenusID, ok := lookups.subgenus[key]; ok {
			relations = append(relations, coldp.NameRelation{
				NameID:                 combination.ID,
				RelatedNameID:          subgenusID,
				Type:                   coldp.NewNomRelType(""),
				TwNameRelationshipType: "TaxonNameRelationship::OriginalCombination::OriginalSubgenus",
			})
		}
	}

	// OriginalSpecies: create if basionym has a specific epithet
	// Skip if it would point to the same name as the basionym (BASIONYM already covers it)
	if basionym.SpecificEpithet != "" {
		key := basionym.Genus + "|" + basionym.SpecificEpithet
		if speciesID, ok := lookups.species[key]; ok {
			if speciesID != basionym.ID {
				relations = append(relations, coldp.NameRelation{
					NameID:                 combination.ID,
					RelatedNameID:          speciesID,
					Type:                   coldp.NewNomRelType(""),
					TwNameRelationshipType: "TaxonNameRelationship::OriginalCombination::OriginalSpecies",
				})
			}
		}
	}

	// OriginalSubspecies/Variety/Form: create if basionym was at an infraspecific rank
	// Skip if it would point to the same name as the basionym (BASIONYM already covers it)
	if basionym.InfraspecificEpithet != "" {
		origType := inferOriginalInfraspecificType(basionym.Rank)
		rank := strings.ToLower(basionym.Rank)
		key := basionym.Genus + "|" + basionym.SpecificEpithet + "|" +
			basionym.InfraspecificEpithet + "|" + rank
		if infraspecificID, ok := lookups.infraspecific[key]; ok {
			if infraspecificID != basionym.ID {
				relations = append(relations, coldp.NameRelation{
					NameID:                 combination.ID,
					RelatedNameID:          infraspecificID,
					Type:                   coldp.NewNomRelType(""),
					TwNameRelationshipType: origType,
				})
			}
		}
	}

	return relations
}

// inferOriginalTypeFromRank determines the appropriate OriginalX type
// based on the name's rank. This is used for the BASIONYM relation to indicate
// which level the original combination covers.
// For species → OriginalSpecies, subspecies → OriginalSubspecies, etc.
func inferOriginalTypeFromRank(rank string) string {
	rank = strings.ToLower(rank)
	switch {
	case strings.Contains(rank, "subform"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubform"
	case strings.Contains(rank, "form"):
		return "TaxonNameRelationship::OriginalCombination::OriginalForm"
	case strings.Contains(rank, "subvariety"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubvariety"
	case strings.Contains(rank, "variety"):
		return "TaxonNameRelationship::OriginalCombination::OriginalVariety"
	case strings.Contains(rank, "subspecies"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubspecies"
	case strings.Contains(rank, "species"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSpecies"
	default:
		// For ranks above species or unknown, default to OriginalSpecies
		// since most combinations are at species level
		return "TaxonNameRelationship::OriginalCombination::OriginalSpecies"
	}
}

// inferOriginalInfraspecificType determines the appropriate OriginalX type
// based on the basionym's rank (subspecies, variety, form, etc.)
// Used when creating additional OriginalX relationships beyond the BASIONYM.
func inferOriginalInfraspecificType(rank string) string {
	rank = strings.ToLower(rank)
	switch {
	case strings.Contains(rank, "subform"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubform"
	case strings.Contains(rank, "form"):
		return "TaxonNameRelationship::OriginalCombination::OriginalForm"
	case strings.Contains(rank, "subvariety"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubvariety"
	case strings.Contains(rank, "variety"):
		return "TaxonNameRelationship::OriginalCombination::OriginalVariety"
	case strings.Contains(rank, "subspecies"):
		return "TaxonNameRelationship::OriginalCombination::OriginalSubspecies"
	default:
		// Default to subspecies for unknown infraspecific ranks
		return "TaxonNameRelationship::OriginalCombination::OriginalSubspecies"
	}
}
