package isfga

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (a *isfga) InsertSynonyms(data []coldp.Synonym) error {
	tx, err := a.db.Begin()
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
  INSERT INTO synonym
  (
    col__id, col__taxon_id, col__source_id, col__name_id, col__name_phrase,
    col__according_to_id, col__status_id, col__reference_id,
    col__link, col__remarks, col__modified, col__modified_by
  )
  VALUES (?,?,?,?,?, ?,?,?, ?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err = stmt.Exec(
			d.ID, d.TaxonID, d.SourceID, d.NameID, d.NamePhrase,
			d.AccordingToID, d.Status.ID(), d.ReferenceID,
			d.Link, d.Remarks, d.Modified, d.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
