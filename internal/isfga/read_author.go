package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadAuthors(
	ctx context.Context,
	ch chan<- coldp.Author,
) error {
	q := `
SELECT
	col__id, col__source_id, col__alternative_id, col__given, col__family,
	col__suffix, col__abbreviation_botany, col__alternative_names, col__sex_id,
	col__country, col__birth, col__birth_place, col__death, col__affiliation,
	col__interest, col__reference_id, col__link, col__remarks, col__modified,
	col__modified_by
FROM author
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA authors: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var au coldp.Author

		var sex string
		err = rows.Scan(
			&au.ID, &au.SourceID, &au.AlternativeID, &au.Given, &au.Family, &au.Suffix,
			&au.AbbreviationBotany, &au.AlternativeNames, &sex, &au.Country, &au.Birth,
			&au.BirthPlace, &au.Death, &au.Affiliation, &au.Interest, &au.ReferenceID,
			&au.Link, &au.Remarks, &au.Modified, &au.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		au.Sex = coldp.NewSex(sex)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- au:
			if count%50_000 == 0 {
				util.Progress(count, "author")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating author rows: %w", err)
	}

	return nil
}
