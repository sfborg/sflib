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
	n.col__id, n.col__alternative_id, n.col__source_id, n.col__scientific_name,
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

	t.col__id, t.col__alternative_id, t.col__source_id, t.col__parent_id, t.col__ordinal,
	t.col__branch_length, t.col__name_phrase, t.col__according_to_id,
	t.col__according_to_page, t.col__according_to_page_link, t.col__scrutinizer,
	t.col__scrutinizer_id, t.col__status_id, t.col__extinct,
	t.col__temporal_range_start_id, t.col__temporal_range_end_id,
	t.col__environment_id, t.col__species, t.col__section, t.col__subgenus, t.col__genus,
	t.col__subtribe, t.col__tribe, t.col__subfamily, t.col__family, t.col__superfamily,
	t.col__suborder, t.col__order, t.col__subclass, t.col__class, t.col__subphylum,
	t.col__phylum, t.col__kingdom, t.col__reference_id, t.col__link, t.col__remarks,
	t.col__modified, t.col__modified_by

FROM name n
	JOIN taxon t
	  ON t.col__name_id = n.col__id

UNION ALL

SELECT
	n.col__id, n.col__alternative_id, n.col__source_id, n.col__scientific_name,
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

	s.col__id, '', s.col__source_id, s.col__taxon_id, NULL,
	NULL, s.col__name_phrase, s.col__according_to_id,
	'', '', '',
	'', s.col__status_id, NULL,
	'', '',
	'', '', '', '', '',
	'', '', '', '', '',
	'', '', '', '', '',
	'', '', s.col__reference_id, s.col__link, s.col__remarks, s.col__modified,
	s.col__modified_by

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

		var altid string
		var rank, notho, code, nomStatus, taxStatus, gender, start, end, env string
		err = rows.Scan(
			&n.NameAlternativeID, &altid, &n.SourceID, &n.ScientificName,
			&n.Authorship, &rank, &n.Uninomial, &n.Genus, &n.InfragenericEpithet,
			&n.SpecificEpithet, &n.InfraspecificEpithet, &n.CultivarEpithet,
			&notho, &n.OriginalSpelling, &n.CombinationAuthorship,
			&n.CombinationAuthorshipID, &n.CombinationExAuthorship,
			&n.CombinationExAuthorshipID, &n.CombinationAuthorshipYear,
			&n.BasionymAuthorship, &n.BasionymAuthorshipID,
			&n.BasionymExAuthorship, &n.BasionymExAuthorshipID,
			&n.BasionymAuthorshipYear, &code, &nomStatus, &n.ReferenceID,
			&n.PublishedInYear, &n.PublishedInPage, &n.PublishedInPageLink,
			&gender, &n.GenderAgreement, &n.Etymology, &n.Link,
			&n.NameRemarks, &n.Modified, &n.ModifiedBy,

			&n.ID, &n.AlternativeID, &n.SourceID, &n.ParentID, &n.Ordinal,
			&n.BranchLength, &n.NamePhrase, &n.AccordingToID,
			&n.AccordingToPage, &n.AccordingToPageLink, &n.Scrutinizer,
			&n.ScrutinizerID, &taxStatus, &n.Extinct, &start, &end, &env, &n.Species,
			&n.Section, &n.Subgenus, &n.Genus, &n.Subtribe, &n.Tribe,
			&n.Subfamily, &n.Family, &n.Superfamily, &n.Suborder, &n.Order,
			&n.Subclass, &n.Class, &n.Subphylum, &n.Phylum, &n.Kingdom,
			&n.ReferenceID, &n.Link, &n.Remarks, &n.Modified, &n.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		n.NameAlternativeID = "col:" + n.NameAlternativeID
		if altid != "" {
			n.NameAlternativeID += "," + altid
		}

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
