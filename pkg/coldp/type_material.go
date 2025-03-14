package coldp

import "database/sql"

// Type material designated to names. Type material should only be
// associated with the original name, not with a recombination.
type TypeMaterial struct {
	// ID is optional
	ID string

	// SourceID where type material is mentioned. It comes from metadata.
	SourceID string

	// NameID is the id of a name that attached to this type.
	NameID string

	// Citation is a reference where type specimen is attached to name.
	Citation string

	// Status of the type specimen.
	Status TypeStatus

	// InstitutionCode is the name or acronym in use by the institution having
	// custody of the material.
	InstitutionCode string

	// The identifier for the specimen in a collection.
	CatalogNumber string

	// ReferenceID indicating the publication of the type designation.
	ReferenceID string

	// Locality of the type. Ideally from largest area to smallest.
	Locality string

	// Country of the type locality. Preferably as ISO codes.
	Country string

	// Latitude is a decimal latitude of the type locality given in WGS84.
	Latitude sql.NullFloat64

	// Longitute is a decimal longitude of the type locality given in WGS84.
	Longitude sql.NullFloat64

	// Altitue is a decimal longitude of the type locality given in WGS84
	Altitude sql.NullInt64

	// Host is the host organism from which the type specimen was obtained
	// (symbiotype).
	Host string

	// Sex of the type specimen.
	Sex Sex

	// Date the type material was gathered. Recommended to be given as ISO 8601
	// dates.
	Date string

	// Collector of the type specimen.
	Collector string

	// AssociatedSequences to the specimen.
	AssociatedSequences string

	// Link to further information about the specimen, e.g. as provided by
	// the institute holding the collection.
	Link string

	// Remarks about the specimen.
	Remarks string

	// Modified is a timestamp.
	Modified string

	// ModifiedBy a person who last modified the record.
	ModifiedBy string
}

func (t TypeMaterial) Headers() []string {
	return []string{
		"col:id",
		"col:sourceId",
		"col:nameId",
		"col:citation",
		"col:status",
		"col:institutionCode",
		"col:catalogNumber",
		"col:referenceID",
		"col:locality",
		"col:country",
		"col:latitude",
		"col:longitude",
		"col:altitude",
		"col:host",
		"col:sex",
		"col:date",
		"col:collector",
		"col:associatedSequences",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates the TypeMaterial object from a row of data.
func (t TypeMaterial) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	t.ID = row["id"]
	t.SourceID = row["sourceid"]
	t.NameID = row["nameid"]
	t.Citation = row["citation"]
	t.Status = NewTypeStatus(row["status"])
	t.InstitutionCode = row["institutioncode"]
	t.CatalogNumber = row["catalognumber"]
	t.ReferenceID = row["referenceid"]
	t.Locality = row["locality"]
	t.Country = row["country"]
	t.Latitude = ToFloat(row["latitude"])
	t.Longitude = ToFloat(row["longitude"])
	t.Altitude = ToInt(row["altitude"])
	t.Host = row["host"]
	t.Sex = NewSex(row["sex"])
	t.Date = row["date"]
	t.Collector = row["collector"]
	t.AssociatedSequences = row["associatedsequences"]
	t.Link = row["link"]
	t.Remarks = row["remarks"]
	t.Modified = row["modified"]
	t.ModifiedBy = row["modifiedby"]
	return t, warning
}
