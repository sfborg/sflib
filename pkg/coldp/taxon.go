package coldp

import "database/sql"

// Taxon represents a taxonomic data in the CoLDP.
type Taxon struct {
	// ID is the unique identifier for this taxon.
	ID string

	// AlternativeID has alternative identifiers for this taxon.
	// either URI/URN/URL, or scope:id, separated by ','
	AlternativeID string

	// LocalID corresponds to a local identifier, usually an integer. Some
	// sources use it in their URLs and then it might be used by GN to show
	// outlink URL.
	LocalID string

	// GlobalID corresponds to a globally unique identifier like UUID, LSID, DOI
	// etc. If it is used in the source URL GN might use it.
	GlobalID string

	// SourceID is the identifier of the source from metadata.
	SourceID string

	// ParentID is the identifier of the parent taxon.
	ParentID string

	// Ordinal is the used to sort siblings of the same ParentID.
	Ordinal sql.NullInt64

	// BranchLength is the branch length of this taxon in a phylogenetic tree.
	BranchLength sql.NullInt64

	// NameID is the identifier of the name associated with this taxon.
	NameID string

	// NamePhrase is an optional annotation attached to the name in this
	// context (eg `sensu lato` etc).
	NamePhrase string

	// AccordingToID is ReferenceID of the source that this taxon is based on.
	AccordingToID string

	// AccordingToPage is the page number in the source where this taxon is
	// determined.
	AccordingToPage string

	// AccordingToPageLink is a link to the page where this taxon is determined.
	AccordingToPageLink string

	// Scrutinizer is the name of the person who scrutinized this taxon.
	Scrutinizer string

	// ScrutinizerID is the identifier of the scrutinizer ORCID if available.
	ScrutinizerID string

	// ScrutinizerDate is the date of the scrutiny.
	ScrutinizerDate string

	// Provisional indicates taxon is only provisionaly accepted.
	Provisional sql.NullBool

	// ReferenceID is the comma-separated list of references that support
	// this taxon concept.
	ReferenceID string

	// Extinct indicates whether this taxon is extinct.
	Extinct sql.NullBool

	// TemporalRangeStart is the start of the temporal range of this taxon.
	TemporalRangeStart GeoTime

	// TemporalRangeEnd is the end of the temporal range of this taxon.
	TemporalRangeEnd GeoTime

	// Environment is the environments where this taxon lives. Uses Environment
	// controlled vocabulary (comma-separated).
	Environment Environment

	// Species is the species name within this taxon.
	Species string

	// Section is the section name within this taxon.
	Section string

	// Subgenus is the subgenus name within this taxon.
	Subgenus string

	// Genus is the genus name within this taxon.
	Genus string

	// Subtribe is the subtribe name within this taxon.
	Subtribe string

	// Tribe is the tribe name within this taxon.
	Tribe string

	// Subfamily is the subfamily name within this taxon.
	Subfamily string

	// Family is the family name within this taxon.
	Family string

	// Superfamily is the superfamily name within this taxon.
	Superfamily string

	// Suborder is the suborder name within this taxon.
	Suborder string

	// Order is the order name within this taxon.
	Order string

	// Subclass is the subclass name within this taxon.
	Subclass string

	// Class is the class name within this taxon.
	Class string

	// Subphylum is the subphylum name within this taxon.
	Subphylum string

	// Phylum is the phylum name within this taxon.
	Phylum string

	// Kingdom is the kingdom name within this taxon.
	Kingdom string

	// Link is a link to more information about this taxon.
	Link string

	// Remarks are any remarks about this taxon.
	Remarks string

	// Modified is the date when this taxon was last modified.
	Modified string

	// ModifiedBy is the user who last modified this taxon.
	ModifiedBy string
}

func (t Taxon) Headers() []string {
	return []string{
		"col:id",
		"col:alternativeId",
		"col:sourceId",
		"col:parentId",
		"col:ordinal",
		"col:branchLength",
		"col:nameId",
		"col:namePhrase",
		"col:accordingToId",
		"col:accordingToPage",
		"col:accordingToPageLink",
		"col:scrutinizer",
		"col:scrutinizerId",
		"col:provisional",
		"col:extinct",
		"col:temporalRangeStart",
		"col:temporalRangeEnd",
		"col:environment",
		"col:species",
		"col:section",
		"col:subgenus",
		"col:genus",
		"col:subtribe",
		"col:tribe",
		"col:subfamily",
		"col:family",
		"col:superfamily",
		"col:suborder",
		"col:order",
		"col:subclass",
		"col:class",
		"col:subphylum",
		"col:phylum",
		"col:kingdom",
		"col:referenceId",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

// Load populates the Taxon object from a row of data.
func (t Taxon) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	t.ID = row["id"]
	t.AlternativeID = row["alternativeid"]
	t.SourceID = row["sourceid"]
	t.ParentID = row["parentid"]
	t.Ordinal = ToInt(row["ordinal"])
	t.BranchLength = ToInt(row["branchlength"])
	t.NameID = row["nameid"]
	t.NamePhrase = row["namephrase"]
	t.AccordingToID = row["accordingtoid"]
	t.AccordingToPage = row["accordingtopage"]
	t.AccordingToPageLink = row["accordingtopagelink"]
	t.Scrutinizer = row["scrutinizer"]
	t.ScrutinizerID = row["scrutinizerid"]
	t.Provisional = ToBool(row["provisional"])
	t.Extinct = ToBool(row["extinct"])
	t.TemporalRangeStart = NewGeoTime(row["temporalrangestart"])
	t.TemporalRangeEnd = NewGeoTime(row["temporalrangeend"])
	t.Environment = NewEnvironment(row["environment"])
	t.Species = row["species"]
	t.Section = row["section"]
	t.Subgenus = row["subgenus"]
	t.Genus = row["genus"]
	t.Subtribe = row["subtribe"]
	t.Tribe = row["tribe"]
	t.Subfamily = row["subfamily"]
	t.Family = row["family"]
	t.Superfamily = row["superfamily"]
	t.Suborder = row["suborder"]
	t.Order = row["order"]
	t.Subclass = row["subclass"]
	t.Class = row["class"]
	t.Subphylum = row["subphylum"]
	t.Phylum = row["phylum"]
	t.Kingdom = row["kingdom"]
	t.ReferenceID = row["referenceid"]
	t.Link = row["link"]
	t.Remarks = row["remarks"]
	t.Modified = row["modified"]
	t.ModifiedBy = row["modifiedby"]
	return t, warning
}
