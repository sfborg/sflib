package isfga

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/sfga"
)

// nameForInference holds the data needed for basionym matching.
type nameForInference struct {
	ID                   string
	CanonicalStemmed     string
	OriginalAuthorship   string // author without parentheses (basionym author)
	CombinationAuthorship string // author in parentheses
	Genus                string
	Rank                 string
}

// InferBasionyms detects and creates basionym relationships by matching
// stemmed epithets and original authorship across all names in the archive.
//
// Algorithm:
// 1. Build lookup map: (stemmed_epithet + original_authorship) -> name_id
// 2. Names without combination authorship are potential basionyms
// 3. Names with combination authorship are combinations
// 4. Match combinations to basionyms; skip ambiguous matches (multiple candidates)
// 5. Create name_relation records for matches
func (a *isfga) InferBasionyms(ctx context.Context, cfg sfga.BasionymInferenceConfig) error {
	// Check if we should skip due to existing relations
	if cfg.SkipIfRelationsExist {
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
		if n.CombinationAuthorship != "" {
			continue
		}
		if n.OriginalAuthorship == "" {
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

	// Match combinations to basionyms
	var relations []coldp.NameRelation
	matchCount := 0

	for _, n := range names {
		// Only names with combination authorship are combinations
		if n.CombinationAuthorship == "" {
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
		relations = append(relations, coldp.NameRelation{
			NameID:                 n.ID,
			RelatedNameID:          basionym.ID,
			Type:                   coldp.NewNomRelType("BASIONYM"),
			TwNameRelationshipType: "", // Not TW-specific
		})

		// Optionally create OriginalCombination relationships
		if cfg.CreateOriginalCombinations {
			origRels := a.createOriginalCombinationRelations(n, basionym)
			relations = append(relations, origRels...)
		}
	}

	slog.Info("Matched combinations to basionyms", "matches", matchCount)

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

// getBasionymKey creates a lookup key from stemmed epithet + original authorship.
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

	return lowestEpithet + "_" + authorship
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
			COALESCE(col__basionym_authorship, ''),
			COALESCE(col__combination_authorship, ''),
			COALESCE(col__genus, ''),
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
			&n.OriginalAuthorship,
			&n.CombinationAuthorship,
			&n.Genus,
			&n.Rank,
		)
		if err != nil {
			return nil, err
		}
		names = append(names, n)
	}

	return names, rows.Err()
}

// createOriginalCombinationRelations creates OriginalGenus, OriginalSpecies, etc.
// relationships based on what changed between the combination and basionym.
func (a *isfga) createOriginalCombinationRelations(
	combination, basionym nameForInference,
) []coldp.NameRelation {
	var relations []coldp.NameRelation

	// OriginalGenus: if genera differ, basionym provides original genus
	if combination.Genus != "" && basionym.Genus != "" &&
		combination.Genus != basionym.Genus {
		relations = append(relations, coldp.NameRelation{
			NameID:                 basionym.ID, // The basionym IS the original
			RelatedNameID:          combination.ID,
			Type:                   coldp.NewNomRelType(""), // No COLDP equivalent
			TwNameRelationshipType: "TaxonNameRelationship::OriginalCombination::OriginalGenus",
		})
	}

	// OriginalSpecies: the basionym itself provides the original species
	// This is typically self-referential in TaxonWorks, pointing to the protonym
	relations = append(relations, coldp.NameRelation{
		NameID:                 basionym.ID,
		RelatedNameID:          basionym.ID, // Self-referential
		Type:                   coldp.NewNomRelType(""),
		TwNameRelationshipType: "TaxonNameRelationship::OriginalCombination::OriginalSpecies",
	})

	return relations
}
