package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadDistributions(
	ctx context.Context,
	ch chan<- coldp.Distribution,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__area, col__area_id, col__gazetteer_id,
	col__status_id, col__reference_id, col__remarks, col__modified,
	col__modified_by
FROM distribution
`
	rows, err := a.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var dist coldp.Distribution

		var gztr, status string
		err = rows.Scan(
			&dist.TaxonID, &dist.SourceID, &dist.Area, &dist.AreaID,
			&gztr, &status, &dist.ReferenceID, &dist.Remarks,
			&dist.Modified, &dist.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan distribution row: %w", err)
		}
		dist.Gazetteer = coldp.NewGazetteerEnt(gztr)
		dist.Status = coldp.NewDistrStatus(status)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- dist:
			if count%50_000 == 0 {
				util.Progress(count, "distribution")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating distribution rows: %w", err)
	}

	return nil
}
