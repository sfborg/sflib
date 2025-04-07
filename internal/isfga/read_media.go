package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadMedia(
	ctx context.Context,
	ch chan<- coldp.Media,
) error {
	q := `
SELECT
	col__taxon_id, col__source_id, col__url, col__type, col__format,
	col__title, col__created, col__creator, col__license, col__link,
	col__remarks, col__modified, col__modified_by
FROM media
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA medias: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var md coldp.Media

		err = rows.Scan(
			&md.TaxonID, &md.SourceID, &md.URL, &md.Type, &md.Format,
			&md.Title, &md.Created, &md.Creator, &md.License, &md.Link,
			&md.Remarks, &md.Modified, &md.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan media row: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- md:
			if count%50_000 == 0 {
				util.Progress(count, "media")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating media rows: %w", err)
	}

	return nil
}
