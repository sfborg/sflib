package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertMedia(data []coldp.Media) error {
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
  INSERT INTO media
    (
    col__taxon_id, col__source_id, col__url, col__type, col__format,
    col__title, col__created, col__creator, col__license, col__link,
    col__remarks, col__modified, col__modified_by
    )
  VALUES (?,?,?,?,?, ?,?,?,?,?, ?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.TaxonID, n.SourceID, n.URL, n.Type, n.Format, n.Title, n.Created,
			n.Creator, n.License, n.Link, n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
