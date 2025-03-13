package isfga

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (a *isfga) InsertVernaculars(data []coldp.Vernacular) error {
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
  INSERT INTO vernacular
    (
    col__taxon_id, col__source_id, col__name, col__transliteration,
    col__language, col__preferred, col__country, col__area, col__sex_id,
    col__reference_id, col__remarks, col__modified, col__modified_by
    )
  VALUES (?,?,?,?,?,?, ?,?,?,?,?,?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err = stmt.Exec(
			d.TaxonID, d.SourceID, d.Name, d.Transliteration, d.Language, d.Preferred,
			d.Country, d.Area, d.Sex.ID(), d.ReferenceID, d.Remarks, d.Modified,
			d.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
