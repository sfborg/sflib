package isfga

import (
	"log/slog"

	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) InsertTaxonProperties(data []coldp.TaxonProperty) error {
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
  INSERT INTO taxon_property
    (
    col__taxon_id, col__source_id, col__property, col__value,
    col__reference_id, col__page, col__ordinal, col__remarks, col__modified,
    col__modified_by
    )
  VALUES (?,?,?,?,?,?, ?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.TaxonID, n.SourceID, n.Property, n.Value, n.ReferenceID, n.Page,
			n.Ordinal, n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
