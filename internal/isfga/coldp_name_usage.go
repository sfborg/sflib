package isfga

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (a *isfga) InsertNameUsages(data []coldp.NameUsage) error {
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

	tStmt, err := tx.Prepare(`
  INSERT INTO taxon
    (
    col__id, col__alternative_id, col__source_id, col__parent_id, col__ordinal,
    col__branch_length, col__name_id, col__name_phrase, col__according_to_id,
    col__according_to_page, col__according_to_page_link, col__scrutinizer,
    col__scrutinizer_id, col__scrutinizer_date, col__status_id,
    col__reference_id, col__extinct, col__temporal_range_start_id,
    col__temporal_range_end_id, col__environment_id, col__species,
    col__section, col__subgenus, col__genus, col__subtribe, col__tribe,
    col__subfamily, col__family, col__superfamily, col__suborder, col__order,
    col__subclass, col__class, col__subphylum, col__phylum, col__kingdom,
    col__link, col__remarks, col__modified, col__modified_by
    )
  VALUES (
    ?,?,?,?,?, ?,?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?, ?,?,?,?,?,
    ?,?,?,?,?, ?,?,?,?,?, ?,?,?,?
    )
`)
	if err != nil {
		return err
	}
	defer tStmt.Close()

	nStmt, err := tx.Prepare(`
  INSERT INTO name
    (
    col__id, col__alternative_id, col__source_id, col__scientific_name,
    col__authorship, col__rank_id, col__uninomial, col__genus,
    col__infrageneric_epithet, col__specific_epithet,
    col__infraspecific_epithet, col__cultivar_epithet, col__notho_id,
    col__original_spelling, col__combination_authorship,
    col__combination_authorship_id, col__combination_ex_authorship,
    col__combination_ex_authorship_id, col__combination_authorship_year,
    col__basionym_authorship, col__basionym_authorship_id,
    col__basionym_ex_authorship, col__basionym_ex_authorship_id,
    col__basionym_authorship_year, col__code_id, col__status_id,
    col__reference_id, col__published_in_year, col__published_in_page,
    col__published_in_page_link, col__gender_id, col__gender_agreement,
    col__etymology, col__link, col__remarks, col__modified, col__modified_by,
    gn__scientific_name_string, gn__parse_quality,
		gn__canonical_simple, gn__canonical_full, gn__canonical_stemmed,
		gn__cardinality, gn__virus, gn__hybrid, gn__surrogate, gn__authors,
		gn__id)
  VALUES (?,?,?,?, ?,?,?,?, ?,?, ?,?,?, ?,?, ?,?, ?,?, ?,?, ?,?, ?,?,?, ?,?,?,
    ?,?,?, ?,?,?,?,?, ?,?,?,?,?,?,?,?,?,?,?) 
`)
	if err != nil {
		return err
	}
	defer nStmt.Close()

	sStmt, err := tx.Prepare(`
  INSERT INTO synonym
    (
    col__id, col__taxon_id, col__source_id, col__name_id, col__name_phrase,
    col__according_to_id, col__status_id, col__reference_id, col__link,
    col__remarks, col__modified, col__modified_by
    )
  VALUES (?,?,?,?,?, ?,?,?,?, ?,?,?)`)
	if err != nil {
		return err
	}
	defer sStmt.Close()

	basStmt, err := tx.Prepare(`
  INSERT INTO name_relation
    (col__name_id, col__related_name_id, col__type_id)
  VALUES (?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer basStmt.Close()

	for _, d := range data {
		switch d.TaxonomicStatus {
		case coldp.AcceptedTS, coldp.ProvisionallyAcceptedTS:
			_, err = tStmt.Exec(
				d.ID, d.AlternativeID, d.SourceID, d.ParentID, d.Ordinal, d.BranchLength,
				d.ID, d.NamePhrase, d.AccordingToID, d.AccordingToPage,
				d.AccordingToPageLink, d.Scrutinizer, d.ScrutinizerID,
				d.ScrutinizerDate, d.TaxonomicStatus.ID(), d.ReferenceID, d.Extinct,
				d.TemporalRangeStart.ID(), d.TemporalRangeEnd.ID(),
				d.Environment, d.Species, d.Section, d.Subgenus, d.Genus, d.Subtribe,
				d.Tribe, d.Subfamily, d.Family, d.Superfamily, d.Suborder, d.Order,
				d.Subclass, d.Class, d.Subphylum, d.Phylum, d.Kingdom,
				d.Link, d.Remarks, d.Modified, d.ModifiedBy,
			)
			if err != nil {
				return err
			}
		case coldp.UnknownTaxSt:
			if d.ParentID != "" {
				_, err = tStmt.Exec(
					d.ID, d.AlternativeID, d.SourceID, d.ParentID, d.Ordinal, d.BranchLength,
					d.ID, d.NamePhrase, d.AccordingToID, d.AccordingToPage,
					d.AccordingToPageLink, d.Scrutinizer, d.ScrutinizerID,
					d.ScrutinizerDate, d.TaxonomicStatus.ID(), d.ReferenceID, d.Extinct,
					d.TemporalRangeStart.ID(), d.TemporalRangeEnd.ID(),
					d.Environment, d.Species, d.Section, d.Subgenus, d.Genus, d.Subtribe,
					d.Tribe, d.Subfamily, d.Family, d.Superfamily, d.Suborder, d.Order,
					d.Subclass, d.Class, d.Subphylum, d.Phylum, d.Kingdom,
					d.Link, d.Remarks, d.Modified, d.ModifiedBy,
				)
				if err != nil {
					return err
				}
			}
		default:
			_, err = sStmt.Exec(
				d.ID, d.ParentID, d.SourceID, d.ID, d.NamePhrase, d.AccordingToID,
				d.TaxonomicStatus.ID(), d.ReferenceID, d.Link, d.NameRemarks,
				d.Modified, d.ModifiedBy,
			)
			if err != nil {
				return err
			}
		}

		_, err = nStmt.Exec(
			d.ID, d.NameAlternativeID, d.SourceID, d.ScientificName, d.Authorship,
			d.Rank.ID(), d.Uninomial, d.GenericName, d.InfragenericEpithet,
			d.SpecificEpithet, d.InfraspecificEpithet, d.CultivarEpithet,
			d.Notho.ID(), d.OriginalSpelling, d.CombinationAuthorship,
			d.CombinationAuthorshipID, d.CombinationExAuthorship,
			d.CombinationExAuthorshipID, d.CombinationAuthorshipYear,
			d.BasionymAuthorship, d.BasionymAuthorshipID, d.BasionymExAuthorship,
			d.BasionymExAuthorshipID, d.BasionymAuthorshipYear, d.Code.ID(),
			d.NameStatus.ID(), d.NameReferenceID, d.PublishedInYear,
			d.PublishedInPage, d.PublishedInPageLink, d.Gender.ID(),
			d.GenderAgreement, d.Etymology, d.Link, d.NameRemarks, d.Modified,
			d.ModifiedBy, d.ScientificNameString, d.ParseQuality, d.CanonicalSimple,
			d.CanonicalFull, d.CanonicalStemmed, d.Cardinality, d.Virus,
			d.Hybrid, d.Surrogate, d.Authors, d.GnID,
		)

		if d.BasionymID == "" {
			continue
		}
		basStmt.Exec(
			d.ID, d.BasionymID, coldp.Basionym.ID(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
