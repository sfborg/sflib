package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadNames(
	ctx context.Context,
	ch chan<- coldp.Name,
) error {
	q := `
SELECT
	col__id, col__alternative_id, col__source_id, col__scientific_name,
	col__authorship, col__rank_id, col__uninomial, col__genus,
	col__infrageneric_epithet, col__specific_epithet,
	col__infraspecific_epithet, col__cultivar_epithet, col__notho_id,
	col__original_spelling, col__combination_authorship,
	col__combination_authorship_id, col__combination_ex_authorship,
	col__combination_ex_authorship_id, col__combination_authorship_year,
	col__basionym_authorship, col__basionym_authorship_id,
	col__basionym_ex_authorship, col__basionym_ex_authorship_id,
	col__basionym_authorship_year, col__code_id, col__status_id,
	col__reference_id, col__published_in_year, col__published_in_page,
	col__published_in_page_link, col__gender_id, col__gender_agreement,
	col__etymology, col__link, col__remarks, col__modified, col__modified_by
FROM name
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA names: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var n coldp.Name

		var rank, notho, code, status, gender string
		err = rows.Scan(
			&n.ID, &n.AlternativeID, &n.SourceID, &n.ScientificName,
			&n.Authorship, &rank, &n.Uninomial, &n.Genus, &n.InfragenericEpithet,
			&n.SpecificEpithet, &n.InfraspecificEpithet, &n.CultivarEpithet,
			&notho, &n.OriginalSpelling, &n.CombinationAuthorship,
			&n.CombinationAuthorshipID, &n.CombinationExAuthorship,
			&n.CombinationExAuthorshipID, &n.CombinationAuthorshipYear,
			&n.BasionymAuthorship, &n.BasionymAuthorshipID,
			&n.BasionymExAuthorship, &n.BasionymExAuthorshipID,
			&n.BasionymAuthorshipYear, &code, &status, &n.ReferenceID,
			&n.PublishedInYear, &n.PublishedInPage, &n.PublishedInPageLink,
			&gender, &n.GenderAgreement, &n.Etymology, &n.Link,
			&n.Remarks, &n.Modified, &n.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		n.Rank = coldp.NewRank(rank)
		n.Notho = coldp.NewNamePart(notho)
		n.Code = coldp.NewNomCode(code)
		n.Status = coldp.NewNomStatus(status)
		n.Gender = coldp.NewGender(gender)

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
