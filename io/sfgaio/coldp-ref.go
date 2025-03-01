package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertReferences(data []coldp.Reference) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			slog.Error("Cannot finish transaction", "error", err)
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
	INSERT INTO reference
		(	
		col__id, col__alternative_id, col__source_id, col__citation, col__type_id,
		col__author, col__author_id, col__editor, col__editor_id, col__title,
		col__title_short, col__container_author, col__container_title,
		col__container_title_short, col__issued, col__accessed,
		col__collection_title, col__collection_editor, col__volume, col__issue,
		col__edition, col__page, col__publisher, col__publisher_place,
		col__version, col__isbn, col__issn, col__doi, col__link, col__remarks,
		col__modified, col__modified_by
		)
	VALUES (
		?,?,?,?,?, ?,?,?,?,?,?, ?,?,?, ?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?
	)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err = stmt.Exec(
			d.ID, d.AlternativeID, d.SourceID, d.Citation, d.Type.ID(),
			d.Author, d.AuthorID, d.Editor, d.EditorID, d.Title, d.TitleShort,
			d.ContainerAuthor, d.ContainerTitle, d.ContainerTitleShort,
			d.Issued, d.Accessed, d.CollectionTitle, d.CollectionEditor,
			d.Volume, d.Issue, d.Edition, d.Page, d.Publisher,
			d.PublisherPlace, d.Version, d.ISBN, d.ISSN, d.DOI,
			d.Link, d.Remarks, d.Modified, d.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
