package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadTreatments(
	ctx context.Context,
	ch chan<- coldp.Treatment,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__document, col__format, col__modified,
	col__modified_by
FROM treatment
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA treatments: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var tr coldp.Treatment

		err = rows.Scan(
			&tr.TaxonID, &tr.SourceID, &tr.Document, &tr.Format,
			&tr.Modified, &tr.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan treatment row: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tr:
			if count%50_000 == 0 {
				util.Progress(count, "treatment")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating treatment rows: %w", err)
	}

	return nil
}
