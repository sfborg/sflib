package coldp

type TaxonConceptRelation struct {
	// TaxonID is a reference to a taxon.
	TaxonID string

	// RelateTaxonID is a reference to the second taxon.
	RelatedTaxonID string

	// SourceID is a reference to source from metadata.
	SourceID string

	Type TaxonConceptRelType

	// ReferenceID to description of TaxonConceptRelation.
	ReferenceID string

	// Remarks about the specimen.
	Remarks string

	// Modified timestamp.
	Modified string

	// ModifiedBy contains name of a person who added/updated the record.
	ModifiedBy string
}

func (t TaxonConceptRelation) Headers() []string {
	return []string{
		"col:taxonId",
		"col:relatedTaxonId",
		"col:sourceId",
		"col:type",
		"col:referenceId",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (t TaxonConceptRelation) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	t.TaxonID = row["taxonid"]
	t.RelatedTaxonID = row["relatedtaxonid"]
	t.SourceID = row["sourceid"]
	t.Type = NewTaxonConceptRelType(row["type"])
	t.ReferenceID = row["referenceid"]
	t.Remarks = row["remarks"]
	t.Modified = row["modified"]
	t.ModifiedBy = row["modifiedby"]
	return t, warning
}
