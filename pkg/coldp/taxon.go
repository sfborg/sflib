package coldp

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gnames/gnlib"
)

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
	Environment []Environment

	// Species is the species name within this taxon.
	Species string

	// SpeciesID is the ID of the taxon's species (SF namespace).
	SpeciesID string

	// Section is the section name within this taxon.
	Section string

	// SectionID is the ID of the taxon's section (SF namespace).
	SectionID string

	// Subgenus is the subgenus name within this taxon.
	Subgenus string

	// SubgenusID is the ID of the taxon's subgenus (SF namespace).
	SubgenusID string

	// Genus is the genus name within this taxon.
	Genus string

	// GenusID is the ID of the taxon's genus (SF namespace).
	GenusID string

	// Subtribe is the subtribe name within this taxon.
	Subtribe string

	// SubtribeID is the ID of the taxon's subtrive (SF namespace).
	SubtribeID string

	// Tribe is the tribe name within this taxon.
	Tribe string

	// TribeID is the ID of the taxon's tribe (SF namespace).
	TribeID string

	// Subfamily is the subfamily name within this taxon.
	Subfamily string

	// SubfamilyID is the ID of the taxon's subfamily (SF namespace).
	SubfamilyID string

	// Family is the family name within this taxon.
	Family string

	// FamilyID is the ID of the taxon's family (SF namespace).
	FamilyID string

	// Superfamily is the superfamily name within this taxon.
	Superfamily string

	// SuperfamilyID is the ID of the taxon's superfamily (SF namespace).
	SuperfamilyID string

	// Suborder is the suborder name within this taxon.
	Suborder string

	// SuborderID is the ID of the taxon's suborder (SF namespace).
	SuborderID string

	// Order is the order name within this taxon.
	Order string

	// OrderID is the ID of the taxon's order (SF namespace).
	OrderID string

	// Subclass is the subclass name within this taxon.
	Subclass string

	// SubclassID is the ID of the taxon's subclass (SF namespace).
	SubclassID string

	// Class is the class name within this taxon.
	Class string

	// ClassID is the ID of the taxon's class (SF namespace).
	ClassID string

	// Subphylum is the subphylum name within this taxon.
	Subphylum string

	// SubphylumID is the ID of the taxon's subphylum (SF namespace).
	SubphylumID string

	// Phylum is the phylum name within this taxon.
	Phylum string

	// PhylumID is the ID of the taxon's phylum (SF namespace).
	PhylumID string

	// Kingdom is the kingdom name within this taxon.
	Kingdom string

	// KingdomID is the id of the taxon's kindom (SF namespace).
	KingdomID string

	// Realm is the realm name within this taxon (SF namespace).
	Realm string

	// RealmID is the ID of the taxon's realm (SF namespace).
	RealmID string

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
		"sf:speciesId",
		"col:section",
		"sf:sectionId",
		"col:subgenus",
		"sf:subgenusId",
		"col:genus",
		"sf:genusId",
		"col:subtribe",
		"sf:subtribeId",
		"col:tribe",
		"sf:tribeId",
		"col:subfamily",
		"sf:subfamilyId",
		"col:family",
		"sf:familyId",
		"col:superfamily",
		"sf:superfamilyId",
		"col:suborder",
		"sf:suborderId",
		"col:order",
		"sf:orderId",
		"col:subclass",
		"sf:subclassId",
		"col:class",
		"sf:classId",
		"col:subphylum",
		"sf:subphylumId",
		"col:phylum",
		"sf:phylumId",
		"col:kingdom",
		"sf:kingdomId",
		"sf:realm",
		"sf:realmId",
		"col:referenceId",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (t Taxon) Row() []string {
	var ordinal, brLen, prov, extinct string
	if t.Ordinal.Valid {
		ordinal = strconv.Itoa(int(t.Ordinal.Int64))
	}
	if t.BranchLength.Valid {
		brLen = strconv.Itoa(int(t.BranchLength.Int64))
	}
	if t.Provisional.Valid {
		prov = strconv.FormatBool(t.Provisional.Bool)
	}
	if t.Extinct.Valid {
		extinct = strconv.FormatBool(t.Extinct.Bool)
	}

	envs := gnlib.Map(t.Environment, func(e Environment) string {
		return e.String()
	})
	envs = gnlib.FilterFunc(envs, func(s string) bool {
		return s != ""
	})

	res := []string{
		t.ID, t.AlternativeID, t.SourceID, t.ParentID, ordinal, brLen,
		t.NameID, t.NamePhrase, t.AccordingToID, t.AccordingToPage,
		t.AccordingToPageLink, t.Scrutinizer, t.ScrutinizerID, prov, extinct,
		t.TemporalRangeStart.String(), t.TemporalRangeEnd.String(),
		strings.Join(envs, ","), t.Species, t.SpeciesID, t.Section, t.SectionID,
		t.Subgenus, t.SubgenusID, t.Genus, t.GenusID, t.Subtribe,
		t.SubtribeID, t.Tribe, t.TribeID, t.Subfamily, t.SubfamilyID,
		t.Family, t.FamilyID, t.Superfamily, t.SuperfamilyID,
		t.Suborder, t.SuborderID, t.Order, t.OrderID, t.Subclass,
		t.SubclassID, t.Class, t.ClassID, t.Subphylum, t.SubphylumID,
		t.Phylum, t.PhylumID, t.Kingdom, t.KingdomID, t.Realm, t.RealmID,
		t.ReferenceID, t.Link, t.Remarks, t.Modified,
		t.ModifiedBy,
	}
	return res
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
	t.Environment = GetEnvironments(row["environment"])
	t.Species = row["species"]
	t.SpeciesID = row["speciesid"]
	t.Section = row["section"]
	t.SectionID = row["sectionid"]
	t.Subgenus = row["subgenus"]
	t.SubgenusID = row["subgenusid"]
	t.Genus = row["genus"]
	t.GenusID = row["genusid"]
	t.Subtribe = row["subtribe"]
	t.SubtribeID = row["subtribeid"]
	t.Tribe = row["tribe"]
	t.TribeID = row["tribeid"]
	t.Subfamily = row["subfamily"]
	t.SubfamilyID = row["subfamilyid"]
	t.Family = row["family"]
	t.FamilyID = row["familyid"]
	t.Superfamily = row["superfamily"]
	t.SuperfamilyID = row["superfamilyid"]
	t.Suborder = row["suborder"]
	t.SuborderID = row["suborderid"]
	t.Order = row["order"]
	t.OrderID = row["orderid"]
	t.Subclass = row["subclass"]
	t.SubclassID = row["subclassid"]
	t.Class = row["class"]
	t.ClassID = row["classid"]
	t.Subphylum = row["subphylum"]
	t.SubphylumID = row["subphylumid"]
	t.Phylum = row["phylum"]
	t.PhylumID = row["phylumid"]
	t.Kingdom = row["kingdom"]
	t.KingdomID = row["kingdomid"]
	t.Realm = row["realm"]
	t.RealmID = row["realmid"]
	t.ReferenceID = row["referenceid"]
	t.Link = row["link"]
	t.Remarks = row["remarks"]
	t.Modified = row["modified"]
	t.ModifiedBy = row["modifiedby"]
	return t, warning
}
