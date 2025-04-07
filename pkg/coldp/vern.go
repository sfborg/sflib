package coldp

import (
	"database/sql"
	"strconv"
)

// Vernacular name of a taxon.
type Vernacular struct {
	// TaxonID of the corresponding taxon.
	TaxonID string

	// SourceID to the source from metadata.
	SourceID string

	// Name (vernacular).
	Name string

	// Transliteration to latin alphabet.
	Transliteration string

	// Langugage of the name (three-letter ISO 639-3 code).
	Language string

	// Preferred is true when this vernacular names is considered the most
	// used for the language.
	Preferred sql.NullBool

	// Country name (two-letter ISO 3166-2  code)
	Country string

	// Area (optional) describing the geographic use of the vernacular name in
	// free text within the given country.
	Area string

	// Sex (optional) of the organism this vernacular name is restricted to.
	Sex Sex

	// ReferenceID where the name is supported.
	ReferenceID string

	// Remarks for the name.
	Remarks string

	// Modified is a timestamp.
	Modified string

	// ModifiedBy a person who last modified the record.
	ModifiedBy string
}

func (v Vernacular) Headers() []string {
	return []string{
		"col:taxonId",
		"col:sourceId",
		"col:name",
		"col:transliteration",
		"col:language",
		"col:preferred",
		"col:country",
		"col:area",
		"col:sex",
		"col:referenceId",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (v Vernacular) Row() []string {
	var pref string
	if v.Preferred.Valid {
		pref = strconv.FormatBool(v.Preferred.Bool)
	}
	res := []string{
		v.TaxonID, v.SourceID, v.Name, v.Transliteration, v.Language,
		pref, v.Country, v.Area, v.Sex.String(), v.ReferenceID,
		v.Remarks, v.Modified, v.ModifiedBy,
	}
	return res
}

// Load populates a Vernacular object from a row of data.
func (v Vernacular) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	v.TaxonID = row["taxonid"]
	v.SourceID = row["sourceid"]
	v.Name = row["name"]
	v.Transliteration = row["transliteration"]
	v.Language = row["language"]
	v.Preferred = ToBool(row["preferred"])
	v.Country = row["country"]
	v.Area = row["area"]
	v.Sex = NewSex(row["sex"])
	v.ReferenceID = row["referenceid"]
	v.Remarks = row["remarks"]
	v.Modified = row["modified"]
	v.ModifiedBy = row["modifiedby"]
	return v, warning
}
