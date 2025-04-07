package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadTaxonConceptRelations(
	ctx context.Context,
	ch chan<- coldp.TaxonConceptRelation,
) error {
	q := `
SELECT
	col__taxon_id, col__related_taxon_id, col__source_id, col__type_id,
	col__reference_id, col__remarks, col__modified, col__modified_by
FROM taxon_concept_relation
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA taxon concept relations: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var tcr coldp.TaxonConceptRelation

		var typ string
		err = rows.Scan(
			&tcr.TaxonID, &tcr.RelatedTaxonID, &tcr.SourceID, &typ,
			&tcr.ReferenceID, &tcr.Remarks, &tcr.Modified, &tcr.ModifiedBy,
		)
		if err != nil {
			return err
		}

		tcr.Type = coldp.NewTaxonConceptRelType(typ)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tcr:
			if count%50_000 == 0 {
				util.Progress(count, "taxon concept relation")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating taxon concept relation rows: %w", err)
	}

	return nil
}
