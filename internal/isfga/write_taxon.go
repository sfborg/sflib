package isfga

import (
	"log/slog"
	"strings"

	"github.com/gnames/gnlib"
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
    col__id, col__alternative_id,
    gn__local_id, gn__global_id, tw__otu_id,
    col__source_id, col__parent_id,
    col__ordinal, col__branch_length, col__name_id, col__name_phrase,
    col__according_to_id, col__according_to_page,
    col__according_to_page_link, col__scrutinizer, col__scrutinizer_id,
    col__scrutinizer_date, col__status_id, col__reference_id, col__extinct,
    col__temporal_range_start_id, col__temporal_range_end_id,
    col__environment_id, col__species, sf__species_id, col__section,
    sf__section_id, col__subgenus, sf__subgenus_id, col__genus, sf__genus_id,
    col__subtribe, sf__subtribe_id, col__tribe, sf__tribe_id, col__subfamily,
    sf__subfamily_id, col__family, sf__family_id, col__superfamily,
    sf__superfamily_id, col__suborder, sf__suborder_id, col__order,
    sf__order_id, col__subclass, sf__subclass_id, col__class, sf__class_id,
    col__subphylum, sf__subphylum_id, col__phylum, sf__phylum_id, col__kingdom,
    sf__kingdom_id, sf__realm, sf__realm_id,
    col__link, col__remarks, col__modified, col__modified_by
    )
  VALUES (
    ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?,
    ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?,?,?
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
		env := gnlib.Map(t.Environment, func(env coldp.Environment) string {
			return env.ID()
		})
		_, err = stmt.Exec(
			t.ID, t.AlternativeID,
			t.LocalID, t.GlobalID, t.OtuID,
			t.SourceID, t.ParentID, t.Ordinal, t.BranchLength,
			t.NameID, t.NamePhrase, t.AccordingToID, t.AccordingToPage,
			t.AccordingToPageLink, t.Scrutinizer, t.ScrutinizerID,
			t.ScrutinizerDate, status.ID(), t.ReferenceID, t.Extinct,
			t.TemporalRangeStart.ID(), t.TemporalRangeEnd.ID(),
			strings.Join(env, ","), t.Species, t.SpeciesID, t.Section, t.SectionID,
			t.Subgenus, t.SubgenusID, t.Genus, t.GenusID, t.Subtribe, t.SubtribeID,
			t.Tribe, t.TribeID, t.Subfamily, t.SubfamilyID, t.Family, t.FamilyID,
			t.Superfamily, t.SuperfamilyID, t.Suborder, t.SuborderID, t.Order,
			t.OrderID, t.Subclass, t.SubclassID, t.Class, t.ClassID, t.Subphylum,
			t.SubphylumID, t.Phylum, t.PhylumID, t.Kingdom, t.KingdomID,
			t.Realm, t.RealmID,
			t.Link, t.Remarks, t.Modified, t.ModifiedBy,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
