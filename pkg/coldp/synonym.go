package coldp

// Synonym represents a taxonomic synonym with associated metadata.
type Synonym struct {
	// ID is an optional unique identifier for the synonym.
	// If provided, it should not conflict with any taxon IDs.
	ID string

	// TaxonID is the unique identifier for the related taxon.
	TaxonID string

	// SourceID is the identifier of the source from metadata.
	SourceID string

	// NameID is the identifier of the name associated with this taxon.
	NameID string

	// NamePhrase is an optional annotation attached to the name in this
	// context (e.g., "sensu lato").
	NamePhrase string

	// AccordingToID is the ReferenceID of the source that this synonym is based on.
	AccordingToID string

	// Status represents the taxonomic status of the synonym.
	Status TaxonomicStatus

	// ReferenceID is the identifier of the reference for the synonym data.
	ReferenceID string

	// Link is a URL to further information about the synonym.
	Link string

	// Remarks contains any additional notes about the synonym.
	Remarks string

	// Modified is the timestamp of the last modification.
	Modified string

	// ModifiedBy is the name of the person who last modified the data.
	ModifiedBy string
}

func (s Synonym) Headers() []string {
	return []string{
		"col:id",
		"col:taxonId",
		"col:sourceId",
		"col:nameId",
		"col:namePhrase",
		"col:accordingToId",
		"col:status",
		"col:referenceId",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates a Synonym object from a row of data.
func (s Synonym) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	s.ID = row["id"]
	s.TaxonID = row["taxonid"]
	s.SourceID = row["sourceid"]
	s.NameID = row["nameid"]
	s.NamePhrase = row["namephrase"]
	s.AccordingToID = row["accordingtoid"]
	s.Status = NewTaxonomicStatus(row["status"])
	s.ReferenceID = row["referenceid"]
	s.Link = row["link"]
	s.Remarks = row["remarks"]
	s.Modified = row["modified"]
	s.ModifiedBy = row["modifiedby"]
	return s, warning
}
