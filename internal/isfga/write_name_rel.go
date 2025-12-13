package isfga

import (
	"log/slog"

	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) InsertNameRelations(data []coldp.NameRelation) error {
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
  INSERT INTO name_relation
    (
    col__name_id, col__related_name_id, col__source_id, col__type_id,
    tw__name_relationship_type,
    col__page, col__reference_id, col__remarks,
    col__modified, col__modified_by
    )
  VALUES (?,?,?,?,?,?,?,?,?,?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.NameID, n.RelatedNameID, n.SourceID, n.Type.ID(),
			n.TwNameRelationshipType,
			n.Page, n.ReferenceID, n.Remarks,
			n.Modified, n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
