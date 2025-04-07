package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadSynonyms(
	ctx context.Context,
	ch chan<- coldp.Synonym,
) error {
	q := `
SELECT
	col__id, col__taxon_id, col__source_id, col__name_id, col__name_phrase,
	col__according_to_id, col__status_id, col__reference_id, col__link,
	col__remarks, col__modified, col__modified_by
FROM synonym
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA synonyms: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var syn coldp.Synonym

		var status string
		err = rows.Scan(
			&syn.ID, &syn.TaxonID, &syn.SourceID, &syn.NameID, &syn.NamePhrase,
			&syn.AccordingToID, &status, &syn.ReferenceID, &syn.Link,
			&syn.Remarks, &syn.Modified, &syn.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan synonym row: %w", err)
		}

		syn.Status = coldp.NewTaxonomicStatus(status)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- syn:
			if count%50_000 == 0 {
				util.Progress(count, "synonym")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating synonym rows: %w", err)
	}

	return nil
}
