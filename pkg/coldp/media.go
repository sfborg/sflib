package coldp

// Media attached to a taxon.
type Media struct {
	// TaxonID is a reference to a taxon.
	TaxonID string

	// SourceID is a reference to source from metadata.
	SourceID string

	// URL to the media source.
	URL string

	// Type is MIME type (eg image, image/jpeg, audio)
	Type string

	// Format of the media file (NOT IN DOCS).
	Format string

	// Title of the media item.
	Title string

	// Creation date
	Created string

	// Creator name
	Creator string

	// License of the media file.
	License string

	// URL to the page from which media item came from.
	Link string

	// Remarks about the media.
	Remarks string

	// Modified is a timestamp.
	Modified string

	// ModifiedBy is the person who created the data object.
	ModifiedBy string
}

func (m Media) Headers() []string {
	return []string{
		"col:taxonId",
		"col:sourceId",
		"col:url",
		"col:type",
		"col:format",
		"col:title",
		"col:created",
		"col:creator",
		"col:license",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates the Media object from a row of data.
func (m Media) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	m.TaxonID = row["taxonid"]
	m.SourceID = row["sourceid"]
	m.URL = row["url"]
	m.Type = row["type"]
	m.Format = row["format"]
	m.Title = row["title"]
	m.Created = row["created"]
	m.Creator = row["creator"]
	m.License = row["license"]
	m.Link = row["link"]
	m.Remarks = row["remarks"]
	m.Modified = row["modified"]
	m.ModifiedBy = row["modifiedby"]
	return m, warning
}
