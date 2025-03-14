package isfga

import (
	"log/slog"

	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) InsertSpeciesInteractions(
	data []coldp.SpeciesInteraction,
) error {
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
  INSERT INTO species_interaction
    (
    col__taxon_id, col__related_taxon_id, col__source_id,
    col__related_taxon_scientific_name, col__type_id, col__reference_id,
    col__remarks, col__modified, col__modified_by
    )
  VALUES (?,?,?,?, ?,?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.TaxonID, n.RelatedTaxonID, n.SourceID, n.RelatedTaxonScientificName,
			n.Type.ID(), n.ReferenceID, n.Remarks, n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
