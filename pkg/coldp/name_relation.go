package coldp

// NameRelation describes the relationship between two scientific names.
type NameRelation struct {
	// NameID is the unique identifier of the first scientific name.
	NameID string

	// RelatedNameID is the unique identifier of the related scientific name.
	RelatedNameID string

	// SourceID is the optional identifier of the source of the relationship
	// information.
	SourceID string

	// Type specifies the nature of the relationship between the two names (e.g.,
	// Basyonym, Homotypic).
	Type NomRelType

	// ReferenceID is the identifier of the reference that describes this
	// relationship.
	ReferenceID string

	//The exact single page number where the nomenclatural relation was published
	//in the linked reference. If the value spans multiple pages, the first page
	//should be given.
	Page string

	// Remarks provides additional notes or details about the relationship.
	Remarks string

	// Modified is a timestamp.
	Modified string

	// ModifiedBy a person who last modified the record.
	ModifiedBy string
}

func (n NameRelation) Headers() []string {
	return []string{
		"col:nameId",
		"col:relatedNameId",
		"col:sourceId",
		"col:type",
		"col:referenceId",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates a NameRelation object with data from a parsed row.
func (n NameRelation) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	n.NameID = row["nameid"]
	n.RelatedNameID = row["relatednameid"]
	n.SourceID = row["sourceid"]
	n.Type = NewNomRelType(row["type"])
	n.ReferenceID = row["referenceid"]
	n.Remarks = row["remarks"]
	n.Modified = row["modified"]
	n.ModifiedBy = row["modifiedby"]
	return n, warning
}
