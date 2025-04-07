package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadTaxonProperties(
	ctx context.Context,
	ch chan<- coldp.TaxonProperty,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__property, col__value,
	col__reference_id, col__page, col__ordinal, col__remarks, col__modified,
	col__modified_by
FROM taxon_property
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA taxon properties: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var tp coldp.TaxonProperty

		err = rows.Scan(
			&tp.TaxonID, &tp.SourceID, &tp.Property, &tp.Value, &tp.ReferenceID,
			&tp.Page, &tp.Ordinal, &tp.Remarks, &tp.Modified, &tp.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan SFGA taxon properties: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tp:
			if count%50_000 == 0 {
				util.Progress(count, "taxon property")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating taxon property rows: %w", err)
	}

	return nil
}
