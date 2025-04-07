package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadVernaculars(
	ctx context.Context,
	ch chan<- coldp.Vernacular,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__name, col__transliteration,
	col__language, col__preferred, col__country, col__area, col__sex_id,
	col__reference_id, col__remarks, col__modified, col__modified_by
FROM vernacular
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA authors: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var vrn coldp.Vernacular

		var sex string
		err = rows.Scan(
			&vrn.TaxonID, &vrn.SourceID, &vrn.Name, &vrn.Transliteration,
			&vrn.Language, &vrn.Preferred, &vrn.Country, &vrn.Area,
			&sex, &vrn.ReferenceID, &vrn.Remarks, &vrn.Modified, &vrn.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		vrn.Sex = coldp.NewSex(sex)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- vrn:
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
