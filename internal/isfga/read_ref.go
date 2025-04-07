package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadReferences(
	ctx context.Context,
	ch chan<- coldp.Reference,
) error {
	q := `
SELECT
	col__id, col__alternative_id, col__source_id, col__citation, col__type_id,
	col__author, col__author_id, col__editor, col__editor_id, col__title,
	col__title_short, col__container_author, col__container_title,
	col__container_title_short, col__issued, col__accessed,
	col__collection_title, col__collection_editor, col__volume, col__issue,
	col__edition, col__page, col__publisher, col__publisher_place,
	col__version, col__isbn, col__issn, col__doi, col__link, col__remarks,
	col__modified, col__modified_by
FROM reference
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA reference: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var ref coldp.Reference

		var typ string
		err = rows.Scan(
			&ref.ID, &ref.AlternativeID, &ref.SourceID, &ref.Citation,
			&typ, &ref.Author, &ref.AuthorID, &ref.Editor,
			&ref.EditorID, &ref.Title, &ref.TitleShort, &ref.ContainerAuthor,
			&ref.ContainerTitle, &ref.ContainerTitleShort,
			&ref.Issued, &ref.Accessed, &ref.CollectionTitle,
			&ref.CollectionEditor, &ref.Volume, &ref.Issue, &ref.Edition,
			&ref.Page, &ref.Publisher, &ref.PublisherPlace, &ref.Version,
			&ref.ISBN, &ref.ISSN, &ref.DOI, &ref.Link, &ref.Remarks,
			&ref.Modified, &ref.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan reference row: %w", err)
		}

		ref.Type = coldp.NewReferenceType(typ)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- ref:
			if count%50_000 == 0 {
				util.Progress(count, "reference")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating reference rows: %w", err)
	}

	return nil
}
