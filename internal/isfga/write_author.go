package isfga

import (
	"log/slog"

	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) InsertAuthors(data []coldp.Author) error {
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
	INSERT INTO author
		(
		col__id, col__source_id, col__alternative_id, col__given, col__family,
		col__suffix, col__abbreviation_botany, col__alternative_names, col__sex_id,
		col__country, col__birth, col__birth_place, col__death, col__affiliation,
		col__interest, col__reference_id, col__link, col__remarks, col__modified,
		col__modified_by
		)
	VALUES (?,?,?,?,?, ?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range data {
		_, err = stmt.Exec(
			n.ID, n.SourceID, n.AlternativeID, n.Given, n.Family,
			n.Suffix, n.AbbreviationBotany, n.AlternativeNames, n.Sex,
			n.Country, n.Birth, n.BirthPlace, n.Death, n.Affiliation,
			n.Interest, n.ReferenceID, n.Link, n.Remarks, n.Modified,
			n.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
