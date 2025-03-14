package coldp

import "database/sql"

// SpeciesEstimate provides estimation of how many species are children of
// the taxon.
type SpeciesEstimate struct {
	// TaxonID is a reference to a taxon.
	TaxonID string

	// SourceID is a reference to source from metadata.
	SourceID string

	// Estimate is the estimated number of species for the taxon.
	Estimate sql.NullInt64

	// Type of the estimation.
	Type EstimateType

	// ReferenceID indicating the publication of the type designation.
	ReferenceID string

	// Remarks about the specimen.
	Remarks string

	// Modified timestamp.
	Modified string

	// ModifiedBy contains name of a person who added/updated the record.
	ModifiedBy string
}

func (s SpeciesEstimate) Headers() []string {
	return []string{
		"col:taxonId",
		"col:sourceId",
		"col:estimate",
		"col:type",
		"col:referenceId",
		"col:remarks",
		"col:modified",
		"col:modified_by",
	}
}

func (s SpeciesEstimate) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	s.TaxonID = row["taxonid"]
	s.SourceID = row["sourceid"]
	s.Estimate = ToInt(row["estimate"])
	s.Type = NewEstimateType(row["type"])
	s.ReferenceID = row["referenceid"]
	s.Remarks = row["remarks"]
	s.Modified = row["modified"]
	s.ModifiedBy = row["modifiedby"]
	return s, warning
}
