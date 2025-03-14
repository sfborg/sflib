package isfga

import (
	"log/slog"

	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) InsertTaxa(data []coldp.Taxon) error {
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
  INSERT INTO taxon
    (
    col__id, col__alternative_id, col__source_id, col__parent_id,
    col__ordinal, col__branch_length, col__name_id, col__name_phrase,
    col__according_to_id, col__according_to_page,
    col__according_to_page_link, col__scrutinizer, col__scrutinizer_id,
    col__scrutinizer_date, col__status_id, col__reference_id, col__extinct,
    col__temporal_range_start_id, col__temporal_range_end_id,
    col__environment_id, col__species, col__section, col__subgenus,
    col__genus, col__subtribe, col__tribe, col__subfamily, col__family,
    col__superfamily, col__suborder, col__order, col__subclass, col__class,
    col__subphylum, col__phylum, col__kingdom, col__link, col__remarks,
    col__modified, col__modified_by
    )
  VALUES (
    ?,?,?,?,?,?, ?,?,?,?, ?,?,?, ?,?,?,?, ?,?, ?,?,?,?,?,?,
    ?,?,?,?,?,?, ?,?,?,?,?, ?,?,?,?
    )
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range data {
		status := coldp.AcceptedTS
		if t.Provisional.Bool {
			status = coldp.ProvisionallyAcceptedTS
		}
		_, err = stmt.Exec(
			t.ID, t.AlternativeID, t.SourceID, t.ParentID, t.Ordinal, t.BranchLength,
			t.NameID, t.NamePhrase, t.AccordingToID, t.AccordingToPage,
			t.AccordingToPageLink, t.Scrutinizer, t.ScrutinizerID,
			t.ScrutinizerDate, status.ID(), t.ReferenceID, t.Extinct,
			t.TemporalRangeStart.ID(), t.TemporalRangeEnd.ID(),
			t.Environment, t.Species, t.Section, t.Subgenus, t.Genus, t.Subtribe,
			t.Tribe, t.Subfamily, t.Family, t.Superfamily, t.Suborder, t.Order,
			t.Subclass, t.Class, t.Subphylum, t.Phylum, t.Kingdom,
			t.Link, t.Remarks, t.Modified, t.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
