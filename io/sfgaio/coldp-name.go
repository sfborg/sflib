package sfgaio

import (
	"log/slog"

	"github.com/gnames/coldp/ent/coldp"
)

func (s *sfgaio) InsertNames(data []coldp.Name) error {
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
    col__etymology, col__link, col__remarks, col__modified,
    col__modified_by, gn__scientific_name_string, gn__parse_quality,
		gn__canonical_simple, gn__canonical_full, gn__canonical_stemmed,
		gn__cardinality, gn__virus, gn__hybrid, gn__surrogate, gn__authors,
		gn__id)
  VALUES (?,?,?,?, ?,?,?,?, ?,?, ?,?,?, ?,?, ?,?, ?,?, ?,?, ?,?, ?,?,?, ?,?,?,
    ?,?,?, ?,?,?,?, ?,?,?, ?,?,?, ?,?,?,?,?, ?) 
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	relStmt, err := tx.Prepare(`
	INSERT INTO name_relation
		(name_id, related_name_id, type_id)
	VALUES (?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer relStmt.Close()

	for _, n := range data {

		_, err = stmt.Exec(
			n.ID, n.AlternativeID, n.SourceID, n.ScientificName, n.Authorship,
			n.Rank.ID(), n.Uninomial, n.Genus, n.InfragenericEpithet,
			n.SpecificEpithet, n.InfraspecificEpithet, n.CultivarEpithet,
			n.Notho.ID(), n.OriginalSpelling, n.CombinationAuthorship,
			n.CombinationAuthorshipID, n.CombinationExAuthorship,
			n.CombinationExAuthorshipID, n.CombinationAuthorshipYear,
			n.BasionymAuthorship, n.BasionymAuthorshipID,
			n.BasionymExAuthorship, n.BasionymExAuthorshipID,
			n.BasionymAuthorshipYear, n.Code.ID(), n.Status.ID(), n.ReferenceID,
			n.PublishedInYear, n.PublishedInPage, n.PublishedInPageLink,
			n.Gender.ID(), n.GenderAgreement, n.Etymology,
			n.Link, n.Remarks, n.Modified, n.ModifiedBy,
			n.ScientificNameString, n.ParseQuality, n.CanonicalSimple,
			n.CanonicalFull, n.CanonicalStemmed, n.Cardinality, n.Virus,
			n.Hybrid, n.Surrogate, n.Authors, n.GnID,
		)
		if err != nil {
			return err
		}

		if n.BasionymID == "" {
			continue
		}
		relStmt.Exec(
			n.ID, n.BasionymID, coldp.Basionym.ID(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
