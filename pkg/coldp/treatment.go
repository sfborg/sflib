package coldp

type Treatment struct {
	// TaxonID is a reference to a taxon.
	TaxonID string

	// SourceID is a reference to source from metadata.
	SourceID string

	// Document (probably path in the file system).
	Document string

	// Format of the document (HTML, TXT, XML etc.).
	Format string

	// Modified timestamp.
	Modified string

	// ModifiedBy contains name of a person who added/updated the record.
	ModifiedBy string
}

func (t Treatment) Headers() []string {
	return []string{
		"col:taxonId",
		"col:sourceId",
		"col:document",
		"col:format",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates the Treatment object from a row of data.
func (t Treatment) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	t.TaxonID = row["taxonid"]
	t.SourceID = row["sourceid"]
	t.Document = row["document"]
	t.Format = row["format"]
	t.Modified = row["modified"]
	t.ModifiedBy = row["modifiedby"]
	return t, warning
}
