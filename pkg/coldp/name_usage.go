package coldp

import (
	"database/sql"
	"strings"
)

// NameUsage combines fields of Name, Taxon and Synonym.
type NameUsage struct {
	ID                        string          // t, s
	AlternativeID             string          // t
	NameAlternativeID         string          // n
	LocalID                   string          // t GN
	GlobalID                  string          // t GN
	SourceID                  string          // n,t,s
	ParentID                  string          // t, s
	BasionymID                string          // t
	TaxonomicStatus           TaxonomicStatus // s
	ScientificName            string          // n
	Authorship                string          // n
	ScientificNameString      string          // GN
	ParseQuality              sql.NullInt64   // GN
	CanonicalSimple           string          // GN
	CanonicalFull             string          // GN
	CanonicalStemmed          string          // GN
	Cardinality               sql.NullInt64   // GN
	Virus                     sql.NullBool    // GN
	Hybrid                    string          // GN
	Surrogate                 string          // GN
	Authors                   string          // GN
	GnID                      string          // GN
	Rank                      Rank            // n
	Notho                     NamePart        // n
	OriginalSpelling          sql.NullBool    // n
	Uninomial                 string          // n
	GenericName               string          // n
	InfragenericEpithet       string          // n
	SpecificEpithet           string          // n
	InfraspecificEpithet      string          // n
	CultivarEpithet           string          // n
	CombinationAuthorship     string          // n
	CombinationAuthorshipID   string          // n
	CombinationExAuthorship   string          // n
	CombinationExAuthorshipID string          // n
	CombinationAuthorshipYear string          // n
	BasionymAuthorship        string          // n
	BasionymAuthorshipID      string          // n
	BasionymExAuthorship      string          // n
	BasionymExAuthorshipID    string          // n
	BasionymAuthorshipYear    string          // n
	NamePhrase                string          // t
	NameReferenceID           string          // n
	PublishedInYear           string          // n
	PublishedInPage           string          // n
	PublishedInPageLink       string          // n
	Gender                    Gender          // n
	GenderAgreement           sql.NullBool    // n
	Etymology                 string          // n
	Code                      NomCode         // n
	NameStatus                NomStatus       // n
	AccordingToID             string          // t
	AccordingToPage           string          // t
	AccordingToPageLink       string          // t
	ReferenceID               string          // t
	Scrutinizer               string          // t
	ScrutinizerID             string          // t
	ScrutinizerDate           string          // t
	Extinct                   sql.NullBool    // t
	TemporalRangeStart        GeoTime         // t
	TemporalRangeEnd          GeoTime         // t
	Environment               Environment     // t
	Species                   string          // t
	Section                   string          // t
	Subgenus                  string          // t
	Genus                     string          // t
	Subtribe                  string          // t
	Tribe                     string          // t
	Subfamily                 string          // t
	Family                    string          // t
	Superfamily               string          // t
	Suborder                  string          // t
	Order                     string          // t
	Subclass                  string          // t
	Class                     string          // t
	Subphylum                 string          // t
	Phylum                    string          // t
	Kingdom                   string          // t
	Ordinal                   sql.NullInt64   // t
	BranchLength              sql.NullInt64   // t
	Link                      string          // n, t
	NameRemarks               string          // n
	Remarks                   string          // t
	Modified                  string          // n, t
	ModifiedBy                string          // n, t
}

func (n NameUsage) Headers() []string {
	return []string{
		"col:id",
		"col:alternativeId",
		"col:nameAlternativeId",
		"col:sourceId",
		"col:parentId",
		"col:basionymId",
		"col:taxonomicStatus",
		"col:scientificName",
		"col:authorship",
		"col:rank",
		"col:notho",
		"col:originalSpelling",
		"col:uninomial",
		"col:genericName",
		"col:infragenericEpithet",
		"col:specificEpithet",
		"col:infraspecificEpithet",
		"col:cultivarEpithet",
		"col:combinationAuthorship",
		"col:combinationAuthorshipId",
		"col:combinationExAuthorship",
		"col:combinationExAuthorshipId",
		"col:combinationAuthorshipYear",
		"col:basionymAuthorship",
		"col:basionymAuthorshipId",
		"col:basionymExAuthorship",
		"col:basionymExAuthorshipId",
		"col:basionymAuthorshipYear",
		"col:namePhrase",
		"col:nameReferenceId",
		"col:publishedInYear",
		"col:publishedInPage",
		"col:publishedInPageLink",
		"col:gender",
		"col:genderAgreement",
		"col:etymology",
		"col:code",
		"col:nameStatus",
		"col:accordingToID",
		"col:accordingToPage",
		"col:accordingToPageLink",
		"col:referenceID",
		"col:scrutinizer",
		"col:scrutinizerId",
		"col:scrutinizerDate",
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
		"col:ordinal",
		"col:branchLength",
		"col:link",
		"col:nameRemarks",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (n NameUsage) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	n.ID = row["id"]
	n.AlternativeID = row["alternativeid"]
	n.NameAlternativeID = row["namealternativeid"]
	n.SourceID = row["sourceid"]
	n.ParentID = row["parentid"]
	n.BasionymID = row["basionymid"]
	n.TaxonomicStatus = NewTaxonomicStatus(row["status"])
	n.ScientificName = row["scientificname"]
	n.Authorship = row["authorship"]
	n.ScientificNameString = n.ScientificName
	if n.Authorship != "" && !strings.HasSuffix(n.ScientificName, n.Authorship) {
		n.ScientificNameString += " " + n.Authorship
	}
	n.Rank = NewRank(row["rank"])
	n.Notho = NewNamePart(row["notho"])
	n.OriginalSpelling = ToBool(row["originalspelling"])
	n.Uninomial = row["uninomial"]
	n.GenericName = row["genericname"]
	n.InfragenericEpithet = row["infragenericepithet"]
	n.SpecificEpithet = row["specificepithet"]
	n.InfraspecificEpithet = row["infraspecificepithet"]
	n.CultivarEpithet = row["cultivarepithet"]
	n.CombinationAuthorship = row["combinationauthorship"]
	n.CombinationAuthorshipID = row["combinationauthorshipid"]
	n.CombinationExAuthorship = row["combinationexauthorship"]
	n.CombinationExAuthorshipID = row["combinationexauthorshipid"]
	n.CombinationAuthorshipYear = row["combinationauthorshipyear"]
	n.BasionymAuthorship = row["basionymauthorship"]
	n.BasionymAuthorshipID = row["basionymauthorshipid"]
	n.BasionymExAuthorship = row["basionymexauthorship"]
	n.BasionymExAuthorshipID = row["basionymexauthorshipid"]
	n.BasionymAuthorshipYear = row["basionymauthorshipyear"]
	n.NamePhrase = row["namephrase"]
	n.NameReferenceID = row["namereferenceid"]
	n.PublishedInYear = row["publishedinyear"]
	n.PublishedInPage = row["publishedinpage"]
	n.PublishedInPageLink = row["publishedinpagelink"]
	n.Gender = NewGender(row["gender"])
	n.GenderAgreement = ToBool(row["genderagreement"])
	n.Etymology = row["etymology"]
	n.Code = NewNomCode(row["code"])
	n.NameStatus = NewNomStatus(row["namestatus"])
	n.AccordingToID = row["accordingtoid"]
	n.AccordingToPage = row["accordingtopage"]
	n.AccordingToPageLink = row["accordingtopagelink"]
	n.ReferenceID = row["referenceid"]
	n.Scrutinizer = row["scrutinizer"]
	n.ScrutinizerID = row["scrutinizerid"]
	n.ScrutinizerDate = row["scrutinizerdate"]
	n.Extinct = ToBool(row["extinct"])
	n.TemporalRangeStart = NewGeoTime(row["temporalrangestart"])
	n.TemporalRangeEnd = NewGeoTime(row["temporalrangeend"])
	n.Environment = NewEnvironment(row["environment"])
	n.Species = row["species"]
	n.Section = row["section"]
	n.Subgenus = row["subgenus"]
	n.Genus = row["genus"]
	n.Subtribe = row["subtribe"]
	n.Tribe = row["tribe"]
	n.Subfamily = row["subfamily"]
	n.Family = row["family"]
	n.Superfamily = row["superfamily"]
	n.Suborder = row["suborder"]
	n.Order = row["order"]
	n.Subclass = row["subclass"]
	n.Class = row["class"]
	n.Subphylum = row["subphylum"]
	n.Phylum = row["phylum"]
	n.Kingdom = row["kingdom"]
	n.Ordinal = ToInt(row["ordinal"])
	n.BranchLength = ToInt(row["branchlength"])
	n.Link = row["link"]
	n.NameRemarks = row["nameremarks"]
	n.Remarks = row["remarks"]
	n.Modified = row["modified"]
	n.ModifiedBy = row["modifiedby"]
	return n, warning
}
