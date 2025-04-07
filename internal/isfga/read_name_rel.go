package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadNameRelationships(
	ctx context.Context,
	ch chan<- coldp.NameRelation,
) error {
	q := `
SELECT
	col__name_id, col__related_name_id, col__source_id, col__type_id,
	col__reference_id, col__remarks, col__modified, col__modified_by
FROM name_relation
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA name relationships: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var nr coldp.NameRelation

		var typ string
		err = rows.Scan(
			&nr.NameID, &nr.RelatedNameID, &nr.SourceID, &typ, &nr.ReferenceID,
			&nr.Remarks, &nr.Modified, &nr.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan name relationship row: %w", err)
		}
		nr.Type = coldp.NewNomRelType(typ)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- nr:
			if count%50_000 == 0 {
				util.Progress(count, "name relationship")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating name relationship rows: %w", err)
	}

	return nil
}
