package coldp

type SpeciesInteraction struct {
	// TaxonID is a reference to a taxon.
	TaxonID string

	// RelateTaxonID is a reference to the second taxon.
	RelatedTaxonID string

	// SourceID is a reference to source from metadata.
	SourceID string

	// RelatedTaxonScientificName is the name of the related taxon.
	// TODO: find if it is deprecated.
	RelatedTaxonScientificName string

	// Type of interaction.
	Type SpInteractionType

	// ReferenceID describing species interaction.
	ReferenceID string

	// Remarks about the specimen.
	Remarks string

	// Modified timestamp.
	Modified string

	// ModifiedBy contains name of a person who added/updated the record.
	ModifiedBy string
}

func (s SpeciesInteraction) Headers() []string {
	return []string{
		"col:taxonId",
		"col:relatedTaxonId",
		"col:sourceId",
		"col:relatedTaxonScientificName",
		"col:type",
		"col:referenceId",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (s SpeciesInteraction) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	s.TaxonID = row["taxonid"]
	s.RelatedTaxonID = row["relatedtaxonid"]
	s.SourceID = row["sourceid"]
	s.RelatedTaxonScientificName = row["relatedtaxonscientificname"]
	s.Type = NewSpInteractionType(row["type"])
	s.ReferenceID = row["referenceid"]
	s.Remarks = row["remarks"]
	s.Modified = row["modified"]
	s.ModifiedBy = row["modifiedby"]
	return s, warning
}
