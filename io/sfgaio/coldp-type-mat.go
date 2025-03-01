package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertTypeMaterials(data []coldp.TypeMaterial) error {
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
  INSERT INTO type_material
    (
    col__id, col__source_id, col__name_id, col__citation, col__status_id,
    col__institution_code, col__catalog_number, col__reference_id,
    col__locality, col__country, col__latitude, col__longitude,
    col__altitude, col__host, col__sex_id, col__date, col__collector,
    col__associated_sequences, col__link, col__remarks, col__modified,
    col__modified_by
    )
  VALUES (?,?,?,?,?, ?,?,?,?, ?,?,?,?,?,?, ?,?,?, ?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.ID, n.SourceID, n.NameID, n.Citation, n.Status.ID(),
			n.InstitutionCode, n.CatalogNumber, n.ReferenceID, n.Locality,
			n.Country, n.Latitude, n.Longitude, n.Altitude, n.Host, n.Sex.ID(),
			n.Date, n.Collector, n.AssociatedSequences,
			n.Link, n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
