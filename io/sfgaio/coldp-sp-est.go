package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertSpeciesEstimates(data []coldp.SpeciesEstimate) error {
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
  INSERT INTO species_estimate
    (
    col__taxon_id, col__source_id, col__estimate, col__type_id,
    col__reference_id, col__remarks, col__modified, col__modified_by
    )
  VALUES (?,?,?,?,?, ?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.TaxonID, n.SourceID, n.Estimate, n.Type.ID(), n.ReferenceID,
			n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
