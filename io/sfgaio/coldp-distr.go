package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertDistributions(data []coldp.Distribution) error {
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
	INSERT INTO distribution
    (
    col__taxon_id, col__source_id, col__area, col__area_id,
    col__gazetteer_id, col__status_id, col__reference_id, col__remarks,
    col__modified, col__modified_by
    )
	VALUES (?,?,?,?, ?,?,?,?, ?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.TaxonID, n.SourceID, n.Area, n.AreaID, n.Gazetteer.ID(),
			n.Status.ID(), n.ReferenceID, n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
