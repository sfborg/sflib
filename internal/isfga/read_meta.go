package isfga

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/pkg/coldp"
)

var validTables = func() gnlib.Set[string] {
	res := make(gnlib.Set[string])
	tables := []string{
		"metadata", "source", "author", "reference",
		"name", "taxon", "synonym", "vernacular",
		"name_relation", "type_material", "distribution",
		"media", "treatment", "species_estimate",
		"taxon_property", "species_interaction",
		"taxon_concept_relation", "name_match",
	}
	for i := range tables {
		res.Add(tables[i])
	}
	return res
}()

func (a *isfga) isEmpty(table string) bool {
	if !validTables.Has(table) {
		return true
	}
	q := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s LIMIT 1)", table)
	var exists bool

	err := a.db.QueryRow(q).Scan(&exists)
	if err != nil {
		return true
	}
	if exists {
		return false
	}
	return true
}

func (a *isfga) LoadMeta() (*coldp.Meta, error) {
	if a.isEmpty("metadata") {
		return nil, nil
	}

	q := `
SELECT
	col__doi, col__title, col__alias, col__description, col__issued,
	col__version, col__keywords, col__geographic_scope, col__taxonomic_scope,
	col__temporal_scope, col__confidence, col__completeness, col__license,
	col__url, col__logo, col__label, col__citation, col__private
  FROM metadata LIMIT 1
`
	res := coldp.Meta{}
	row := a.db.QueryRow(q)
	var keywords string
	err := row.Scan(
		&res.DOI, &res.Title, &res.Alias, &res.Description,
		&res.Issued, &res.Version, &keywords,
		&res.GeographicScope, &res.TaxonomicScope, &res.TemporalScope,
		&res.Confidence, &res.Completeness, &res.License, &res.URL,
		&res.Logo, &res.Label, &res.Citation, &res.Private,
	)
	if err != nil {
		return nil, err
	}
	if keywords != "" {
		words := strings.Split(keywords, ",")
		words = gnlib.Map(words, func(word string) string {
			return strings.TrimSpace(word)
		})
		res.Keywords = words
	}
	res.Contact, err = a.getActor("contact")
	if err != nil {
		return nil, err
	}
	res.Publisher, err = a.getActor("publisher")
	if err != nil {
		return nil, err
	}
	res.Editors, err = a.getActors("editor")
	if err != nil {
		return nil, err
	}
	res.Creators, err = a.getActors("creator")
	if err != nil {
		return nil, err
	}
	res.Contributors, err = a.getActors("contributor")
	if err != nil {
		return nil, err
	}
	res.Sources, err = a.getSources()
	return &res, nil
}

func (a *isfga) getActor(table string) (*coldp.Actor, error) {
	q := fmt.Sprintf(`
SELECT
		col__orcid, col__given, col__family, col__rorid, col__organisation,
	  col__email, col__url, col__note, col__city, col__country, col__state
	FROM %s
	LIMIT 1
`, table)
	row := a.db.QueryRow(q)
	var res coldp.Actor

	err := row.Scan(
		&res.Orcid, &res.Given, &res.Family, &res.RorID, &res.Organization,
		&res.Email, &res.URL, &res.Note, &res.City, &res.Country, &res.State,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *isfga) getActors(table string) ([]coldp.Actor, error) {
	var res []coldp.Actor
	q := fmt.Sprintf(`
SELECT
		col__orcid, col__given, col__family, col__rorid, col__organisation,
	  col__email, col__url, col__note, col__city, col__country, col__state
	FROM %s
`, table)
	rows, err := a.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var act coldp.Actor
		err := rows.Scan(
			&act.Orcid, &act.Given, &act.Family, &act.RorID, &act.Organization,
			&act.Email, &act.URL, &act.Note, &act.City, &act.Country, &act.State,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, act)
	}

	return res, nil
}

func (a *isfga) getSources() ([]coldp.Source, error) {
	var res []coldp.Source
	q := `
SELECT
		col__type, col__title, col__authors, col__issued, col__isbn
	FROM source
`
	rows, err := a.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var src coldp.Source
		var authors string
		err = rows.Scan(&src.Type, &src.Title, &authors, &src.Issued, &src.ISBN)
		if err != nil {
			return nil, err
		}
		es := strings.Split(authors, ",|")
		auAry := gnlib.Map(es, func(s string) any {
			s = strings.TrimSpace(s)
			var res any = s
			return res
		})
		src.Authors = auAry
		res = append(res, src)
	}
	return res, nil
}
