package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadSpeciesEstimates(
	ctx context.Context,
	ch chan<- coldp.SpeciesEstimate,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__estimate, col__type_id,
	col__reference_id, col__remarks, col__modified, col__modified_by
FROM species_estimate
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA species estimates: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var se coldp.SpeciesEstimate

		var typ string
		err = rows.Scan(
			&se.TaxonID, &se.SourceID, &se.Estimate, &typ, &se.ReferenceID,
			&se.Remarks, &se.Modified, &se.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan species estimate row: %w", err)
		}

		se.Type = coldp.NewEstimateType(typ)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- se:
			if count%50_000 == 0 {
				util.Progress(count, "species estimate")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating species estimate rows: %w", err)
	}

	return nil
}
