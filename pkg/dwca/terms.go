package dwca

// DwC / DC / GBIF namespace prefixes.
const (
	DwcNS = "http://rs.tdwg.org/dwc/terms/"
	DcNS  = "http://purl.org/dc/terms/"
	GbifNS = "http://rs.gbif.org/terms/1.0/"
)

// Standard row-type URIs.
const (
	TaxonRowType          = DwcNS + "Taxon"
	VernacularNameRowType = GbifNS + "VernacularName"
	DistributionRowType   = GbifNS + "Distribution"
)

// DwcTerm pairs a short CSV header name with its full term URI.
type DwcTerm struct {
	Name string // CSV column header (e.g. "taxonID")
	URI  string // full term URI for meta.xml
}

// CoreTerms defines the ordered fields for the Taxon core file.
// Based on the Catalogue of Life DwCA (col.xml).
var CoreTerms = []DwcTerm{
	{"taxonID", DwcNS + "taxonID"},
	{"parentNameUsageID", DwcNS + "parentNameUsageID"},
	{"acceptedNameUsageID", DwcNS + "acceptedNameUsageID"},
	{"originalNameUsageID", DwcNS + "originalNameUsageID"},
	{"taxonomicStatus", DwcNS + "taxonomicStatus"},
	{"taxonRank", DwcNS + "taxonRank"},
	{"scientificName", DwcNS + "scientificName"},
	{"scientificNameAuthorship", DwcNS + "scientificNameAuthorship"},
	{"genericName", DwcNS + "genericName"},
	{"infragenericEpithet", DwcNS + "infragenericEpithet"},
	{"specificEpithet", DwcNS + "specificEpithet"},
	{"infraspecificEpithet", DwcNS + "infraspecificEpithet"},
	{"cultivarEpithet", DwcNS + "cultivarEpithet"},
	{"nomenclaturalCode", DwcNS + "nomenclaturalCode"},
	{"nomenclaturalStatus", DwcNS + "nomenclaturalStatus"},
	{"kingdom", DwcNS + "kingdom"},
	{"phylum", DwcNS + "phylum"},
	{"class", DwcNS + "class"},
	{"order", DwcNS + "order"},
	{"family", DwcNS + "family"},
	{"genus", DwcNS + "genus"},
	{"subgenus", DwcNS + "subgenus"},
	{"taxonRemarks", DwcNS + "taxonRemarks"},
	{"references", DcNS + "references"},
	{"modified", DcNS + "modified"},
}

// VernacularTerms defines the ordered fields for the VernacularName
// extension file.
var VernacularTerms = []DwcTerm{
	{"taxonID", DwcNS + "taxonID"},
	{"vernacularName", DwcNS + "vernacularName"},
	{"language", DcNS + "language"},
	{"countryCode", DwcNS + "countryCode"},
	{"source", DcNS + "source"},
}

// DistributionTerms defines the ordered fields for the Distribution
// extension file.
var DistributionTerms = []DwcTerm{
	{"taxonID", DwcNS + "taxonID"},
	{"occurrenceStatus", DwcNS + "occurrenceStatus"},
	{"locationID", DwcNS + "locationID"},
	{"locality", DwcNS + "locality"},
	{"countryCode", DwcNS + "countryCode"},
	{"source", DcNS + "source"},
}

// TermHeaders returns just the header names from a slice of DwcTerms.
func TermHeaders(terms []DwcTerm) []string {
	res := make([]string, len(terms))
	for i, t := range terms {
		res[i] = t.Name
	}
	return res
}
