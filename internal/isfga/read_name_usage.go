package isfga

import (
	"context"
	"fmt"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadNameUsages(
	ctx context.Context,
	ch chan<- coldp.NameUsage,
) error {
	q := `
SELECT
	n.col__id, n.col__alternative_id, n.col__source_id,
	n.tw__taxon_name_id,
	n.gn__scientific_name_string, n.gn__parse_quality, n.gn__canonical_simple,
	n.gn__canonical_full, n.gn__canonical_stemmed, n.gn__cardinality,
	n.gn__virus, n.gn__hybrid, n.gn__surrogate, n.gn__authors, n.gn__id,
	n.col__scientific_name,
	n.col__authorship, n.col__rank_id, n.col__uninomial, n.col__genus,
	n.col__infrageneric_epithet, n.col__specific_epithet,
	n.col__infraspecific_epithet, n.col__cultivar_epithet, n.col__notho_id,
	n.col__original_spelling, n.col__combination_authorship,
	n.col__combination_authorship_id, n.col__combination_ex_authorship,
	n.col__combination_ex_authorship_id, n.col__combination_authorship_year,
	n.col__basionym_authorship, n.col__basionym_authorship_id,
	n.col__basionym_ex_authorship, n.col__basionym_ex_authorship_id,
	n.col__basionym_authorship_year, n.col__code_id, n.col__status_id,
	n.col__reference_id, n.col__published_in_year, n.col__published_in_page,
	n.col__published_in_page_link, n.col__gender_id, n.col__gender_agreement,
	n.col__etymology, n.col__link, n.col__remarks, n.col__modified,
	n.col__modified_by,

	t.col__id, t.col__alternative_id,
	t.gn__local_id, t.gn__global_id, t.tw__otu_id,
	t.col__source_id, t.col__parent_id, t.col__ordinal,
	t.col__branch_length, t.col__name_phrase, t.col__according_to_id,
	t.col__according_to_page, t.col__according_to_page_link, t.col__scrutinizer,
	t.col__scrutinizer_id, t.col__scrutinizer_date, t.col__status_id, t.col__extinct,
	t.col__temporal_range_start_id, t.col__temporal_range_end_id,
	t.col__environment_id,
	t.col__species, t.sf__species_id, t.col__section, t.sf__section_id,
	t.col__subgenus, t.sf__subgenus_id, t.col__genus, t.sf__genus_id,
	t.col__subtribe, t.sf__subtribe_id, t.col__tribe, t.sf__tribe_id,
	t.col__subfamily, t.sf__subfamily_id, t.col__family, t.sf__family_id,
	t.col__superfamily, t.sf__superfamily_id, t.col__suborder, t.sf__suborder_id,
	t.col__order, t.sf__order_id, t.col__subclass, t.sf__subclass_id,
	t.col__class, t.sf__class_id, t.col__subphylum, t.sf__subphylum_id,
	t.col__phylum, t.sf__phylum_id, t.col__kingdom, t.sf__kingdom_id,
	t.sf__realm, t.sf__realm_id,
	t.col__reference_id, t.col__link, t.col__remarks, t.col__modified, t.col__modified_by

FROM name n
	JOIN taxon t
	  ON t.col__name_id = n.col__id

UNION ALL

SELECT
	n.col__id, n.col__alternative_id, n.col__source_id,
	n.tw__taxon_name_id,
	n.gn__scientific_name_string, n.gn__parse_quality, n.gn__canonical_simple,
	n.gn__canonical_full, n.gn__canonical_stemmed, n.gn__cardinality,
	n.gn__virus, n.gn__hybrid, n.gn__surrogate, n.gn__authors, n.gn__id,
	n.col__scientific_name,
	n.col__authorship, n.col__rank_id, n.col__uninomial, n.col__genus,
	n.col__infrageneric_epithet, n.col__specific_epithet,
	n.col__infraspecific_epithet, n.col__cultivar_epithet, n.col__notho_id,
	n.col__original_spelling, n.col__combination_authorship,
	n.col__combination_authorship_id, n.col__combination_ex_authorship,
	n.col__combination_ex_authorship_id, n.col__combination_authorship_year,
	n.col__basionym_authorship, n.col__basionym_authorship_id,
	n.col__basionym_ex_authorship, n.col__basionym_ex_authorship_id,
	n.col__basionym_authorship_year, n.col__code_id, n.col__status_id,
	n.col__reference_id, n.col__published_in_year, n.col__published_in_page,
	n.col__published_in_page_link, n.col__gender_id, n.col__gender_agreement,
	n.col__etymology, n.col__link, n.col__remarks, n.col__modified,
	n.col__modified_by,

	s.col__id, '',
	'', '', '',
	s.col__source_id, s.col__taxon_id, NULL,
	NULL, s.col__name_phrase, s.col__according_to_id,
	'', '', '',
	'', '', s.col__status_id, NULL,
	'', '',
	'',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '', '', '',
	'', '',
	s.col__reference_id, s.col__link, s.col__remarks, s.col__modified, s.col__modified_by

	FROM name n
		JOIN synonym s
			ON n.col__id = s.col__name_id
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA names: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var n coldp.NameUsage

		var name_id string
		var rank, notho, code, nomStatus, taxStatus, gender, start, end, env string
		err = rows.Scan(
			&name_id, &n.AlternativeID, &n.SourceID,
			&n.TwTaxonNameID,
			&n.ScientificNameString, &n.ParseQuality, &n.CanonicalSimple,
			&n.CanonicalFull, &n.CanonicalStemmed, &n.Cardinality,
			&n.Virus, &n.Hybrid, &n.Surrogate, &n.Authors, &n.GnID,
			&n.ScientificName,
			&n.Authorship, &rank, &n.Uninomial, &n.GenericName, &n.InfragenericEpithet,
			&n.SpecificEpithet, &n.InfraspecificEpithet, &n.CultivarEpithet,
			&notho, &n.OriginalSpelling, &n.CombinationAuthorship,
			&n.CombinationAuthorshipID, &n.CombinationExAuthorship,
			&n.CombinationExAuthorshipID, &n.CombinationAuthorshipYear,
			&n.BasionymAuthorship, &n.BasionymAuthorshipID,
			&n.BasionymExAuthorship, &n.BasionymExAuthorshipID,
			&n.BasionymAuthorshipYear, &code, &nomStatus, &n.NameReferenceID,
			&n.PublishedInYear, &n.PublishedInPage, &n.PublishedInPageLink,
			&gender, &n.GenderAgreement, &n.Etymology, &n.Link,
			&n.NameRemarks, &n.Modified, &n.ModifiedBy,

			&n.ID, &n.AlternativeID,
			&n.LocalID, &n.GlobalID, &n.OtuID,
			&n.SourceID, &n.ParentID, &n.Ordinal,
			&n.BranchLength, &n.NamePhrase, &n.AccordingToID,
			&n.AccordingToPage, &n.AccordingToPageLink, &n.Scrutinizer,
			&n.ScrutinizerID, &n.ScrutinizerDate, &taxStatus, &n.Extinct, &start, &end, &env,
			&n.Species, &n.SpeciesID, &n.Section, &n.SectionID,
			&n.Subgenus, &n.SubgenusID, &n.Genus, &n.GenusID,
			&n.Subtribe, &n.SubtribeID, &n.Tribe, &n.TribeID,
			&n.Subfamily, &n.SubfamilyID, &n.Family, &n.FamilyID,
			&n.Superfamily, &n.SuperfamilyID, &n.Suborder, &n.SuborderID,
			&n.Order, &n.OrderID, &n.Subclass, &n.SubclassID,
			&n.Class, &n.ClassID, &n.Subphylum, &n.SubphylumID,
			&n.Phylum, &n.PhylumID, &n.Kingdom, &n.KingdomID,
			&n.Realm, &n.RealmID,
			&n.ReferenceID, &n.Link, &n.Remarks, &n.Modified, &n.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		// TODO: decide how to deal with 'real' name.col__id, which now we
		// ignore
		_ = name_id

		n.Rank = coldp.NewRank(rank)
		n.Notho = coldp.NewNamePart(notho)
		n.Code = nomcode.New(code)
		n.NameStatus = coldp.NewNomStatus(nomStatus)
		n.Gender = coldp.NewGender(gender)

		n.TaxonomicStatus = coldp.NewTaxonomicStatus(taxStatus)
		n.TemporalRangeStart = coldp.NewGeoTime(start)
		n.TemporalRangeEnd = coldp.NewGeoTime(end)

		envs := strings.Split(env, ",")
		envs = gnlib.Map(envs, func(s string) string {
			return strings.TrimSpace(s)
		})
		for _, v := range envs {
			coldpEnv := coldp.NewEnvironment(v)
			if coldpEnv != coldp.UnknownEnv {
				n.Environment = append(n.Environment, coldpEnv)
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- n:
			if count%50_000 == 0 {
				util.Progress(count, "name")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating name rows: %w", err)
	}

	return nil
}
