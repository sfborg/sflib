package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadSpeciesInteractions(
	ctx context.Context,
	ch chan<- coldp.SpeciesInteraction,
) error {
	q := `
SELECT
	col__taxon_id, col__related_taxon_id, col__source_id,
	col__related_taxon_scientific_name, col__type_id, col__reference_id,
	col__remarks, col__modified, col__modified_by
FROM species_interaction
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA species interactions: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var si coldp.SpeciesInteraction

		var typ string
		err = rows.Scan(
			&si.TaxonID, &si.RelatedTaxonID, &si.SourceID,
			&si.RelatedTaxonScientificName, &typ, &si.ReferenceID, &si.Remarks,
			&si.Modified, &si.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan species interaction row: %w", err)
		}
		si.Type = coldp.NewSpInteractionType(typ)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- si:
			if count%50_000 == 0 {
				util.Progress(count, "species interaction")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating species interaction rows: %w", err)
	}

	return nil
}
