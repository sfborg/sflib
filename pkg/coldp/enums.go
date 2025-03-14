package coldp

import "strings"

// NomCode provides types of nomenclatural Codes.
type NomCode int

// Constants for different nomenclatural codes.
const (
	UnknownNomCode    NomCode = iota
	Bacterial                 // Bacteriological Code
	Botanical                 // Botanical Code
	Cultivars                 // Cultivated Plant Code
	PhytoSociological         // Phytosociological Code
	Virus                     // Virus Code
	Zoological                // Zoological Code
)

// NewNomCode converts a string (number or word) to NomCode.
func NewNomCode(s string) NomCode {
	s = strings.ToLower(s)
	switch s {
	case "1", "bacterial", "icnp":
		return Bacterial
	case "2", "botanical", "icn", "icnafp", "icbn":
		return Botanical
	case "3", "cultivars", "icncp":
		return Cultivars
	case "4", "phytosociological", "icpn":
		return PhytoSociological
	case "5", "virus", "icvcn":
		return Virus
	case "6", "zoological", "iczn":
		return Zoological
	default:
		return UnknownNomCode
	}
}

var nomCodeToString = map[NomCode]string{
	Bacterial:         "BACTERIAL",
	Botanical:         "BOTANICAL",
	Cultivars:         "CULTIVARS",
	PhytoSociological: "PHYTOSOCIOLOGICAL",
	Virus:             "VIRUS",
	Zoological:        "ZOOLOGICAL",
}

func (nc NomCode) ID() string {
	if res, ok := nomCodeToString[nc]; ok {
		return res
	}
	return ""
}

func (nc NomCode) String() string {
	return ToStr(nc.ID())
}

type NomRelType int

const (
	UnknownNomRelType NomRelType = iota
	SpellingCorrection
	Basionym
	BasedOn
	ReplacementName
	ConservedNRT
	LaterHomonym
	Superfluous
	Homotypic
	Type
)

var nomRelTypeToString = map[NomRelType]string{
	SpellingCorrection: "SPELLING_CORRECTION",
	Basionym:           "BASIONYM",
	BasedOn:            "BASEDON",
	ReplacementName:    "REPLACEMENT_NAME",
	ConservedNRT:       "CONSERVED",
	LaterHomonym:       "LATER_HOMONYM",
	Superfluous:        "SUPERFLUOUS",
	Homotypic:          "HOMOTYPIC",
	Type:               "TYPE",
}

var stringToNomRelType = func() map[string]NomRelType {
	res := make(map[string]NomRelType)
	for k, v := range nomRelTypeToString {
		v = strings.ReplaceAll(v, " ", "")
		v = strings.ReplaceAll(v, "_", "")
		res[v] = k
	}
	return res
}()

func NewNomRelType(s string) NomRelType {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "_", "")
	if res, ok := stringToNomRelType[s]; ok {
		return res
	}
	return UnknownNomRelType
}

func (n NomRelType) ID() string {
	if res, ok := nomRelTypeToString[n]; ok {
		return res
	}
	return ""
}

func (n NomRelType) String() string {
	return ToStr(n.ID())
}

// NamePart represents the part of a scientific name.
type NamePart int

// Constants for different name parts.
const (
	UnknownNP       NamePart = iota
	GenericNP                // The genus part of the name.
	InfragenericNP           // The infrageneric epithet (e.g., subgenus, section) of the name.
	SpecificNP               // The specific epithet of the name.
	InfraspecificNP          // The infraspecific epithet (e.g., subspecies, variety) of the name.
)

// NewNamePart creates a new NamePart from a string representation.
func NewNamePart(s string) NamePart {
	s = strings.ToUpper(s)
	switch s {
	case "GENERIC":
		return GenericNP
	case "INFRAGENERIC":
		return InfragenericNP
	case "SPECIFIC":
		return SpecificNP
	case "INFRASPECIFIC":
		return InfraspecificNP
	default:
		return UnknownNP
	}
}

var namePartToString = map[NamePart]string{
	GenericNP:       "GENERIC",
	InfragenericNP:  "INFRAGENERIC",
	SpecificNP:      "SPECIFIC",
	InfraspecificNP: "INFRASPECIFIC",
}

func (np NamePart) ID() string {
	if res, ok := namePartToString[np]; ok {
		return res
	}
	return ""
}

func (np NamePart) String() string {
	return ToStr(np.ID())
}

// ArchiveType provides type of CoLDP. Can be 'flat' or NameUsage, and
// 'original' or Name.
type ArchiveType int

// Constants for different archive types.
const (
	UnknownAT   ArchiveType = iota
	NameAT                  // Represents an archive with Name data.
	NameUsageAT             // Represents an archive with NameUsage data.
)

// Environment represents the environment type of a taxon.
type Environment int

// Constants for different environment types.
const (
	UnknownEnv  Environment = iota
	Brackish                // Living in brackish water (mix of fresh and salt water).
	Freshwater              // Living in freshwater.
	Marine                  // Living in saltwater.
	Terrestrial             // Living on land.
)

// envToString maps Environment values to their string representations.
var envToString = map[Environment]string{
	Brackish:    "BRACKISH",
	Freshwater:  "FRESHWATER",
	Marine:      "MARINE",
	Terrestrial: "TERRESTRIAL",
}

// stringToEnv maps string representations to Environment values.
var stringToEnv = func() map[string]Environment {
	res := make(map[string]Environment)
	for k, v := range envToString {
		res[v] = k
	}
	return res
}()

// ID returns string ID of the Environment.
func (e Environment) ID() string {
	if res, ok := envToString[e]; ok {
		return res
	}
	return ""
}

// String returns the string representation of the Environment.
func (e Environment) String() string {
	return ToStr(e.ID())
}

// NewEnvironment creates a new Environment from a string representation.
func NewEnvironment(s string) Environment {
	s = strings.ToUpper(s)
	if res, ok := stringToEnv[s]; ok {
		return res
	}
	return UnknownEnv
}

// Sex represents biological sex of a person or organism.
type Sex int

const (
	UnknownSex Sex = iota
	Female
	Male
	Hermaphrodite
)

// NewSex creates a new Sex object from a string
func NewSex(s string) Sex {
	s = strings.ToUpper(s)
	s = strings.TrimRight(s, ".")
	switch s {
	case "M", "MALE":
		return Male
	case "F", "FEMALE":
		return Female
	case "HERM", "HERMAPHRODITE":
		return Hermaphrodite
	default:
		return UnknownSex
	}
}

// ID returns string ID for Sex.
func (s Sex) ID() string {
	switch s {
	case Male:
		return "MALE"
	case Female:
		return "FEMALE"
	case Hermaphrodite:
		return "HERMAPHRODITE"
	default:
		return ""
	}
}

// String returns string representation of Sex object.
func (s Sex) String() string {
	return ToStr(s.ID())
}

// Gender represents the grammatical gender of a scientific name.
type Gender int

// Constants for different grammatical genders.
const (
	UnknownG Gender = iota
	Masculine
	Feminine
	Neutral
)

// NewGender creates a new Gender from a string representation.
func NewGender(s string) Gender {
	s = strings.ToUpper(s)
	s = strings.TrimRight(s, ".")
	switch s {
	case "M", "MASCULINE":
		return Masculine
	case "F", "FEMININE":
		return Feminine
	case "N", "NEUT", "NEUTRAL":
		return Neutral
	default:
		return UnknownG
	}
}

// ID returns string ID for Gender.
func (g Gender) ID() string {
	switch g {
	case Masculine:
		return "MASCULINE"
	case Feminine:
		return "FEMININE"
	case Neutral:
		return "NEUTRAL"
	default:
		return ""
	}
}

// String returns the string representation of the gender.
func (g Gender) String() string {
	return ToStr(g.ID())
}

// FileType represents the type of file.
type FileType int

// Constants for different file types.
const (
	UnknownFileType FileType = iota
	JSON
	YAML
	CSV
	TSV
	// pipe-separated file
	PSV
	// json-based Citation Style Language
	JSONCSL
	BIBTEX
)

// DataType provides types of data files in CoLDP archive.
// It is used to convert a file name to the type of information
// provided in that file.
type DataType int

// Constants for different data types.
const (
	UnkownDT DataType = iota
	AuthorDT
	DistributionDT
	MediaDT
	NameDT
	NameRelationDT
	NameUsageDT
	ReferenceDT
	ReferenceJsonDT
	ReferenceBibtexDT
	SpeciesEstimateDT
	SpeciesInteractionDT
	SynonymDT
	TaxonConceptRelationDT
	TaxonDT
	TaxonPropertyDT
	TreatmentDT
	TypeMaterialDT
	VernacularNameDT
)

// FileFormats returns a list of supported file formats for a given DataType.
func (dt DataType) FileFormats() []FileType {
	switch dt {
	case AuthorDT, NameRelationDT, TaxonDT, SynonymDT, NameUsageDT,
		NameDT, TaxonPropertyDT, TaxonConceptRelationDT, SpeciesInteractionDT,
		SpeciesEstimateDT, TypeMaterialDT, DistributionDT, MediaDT,
		VernacularNameDT, TreatmentDT:
		return []FileType{CSV, PSV, TSV}
	case ReferenceDT:
		return []FileType{CSV, PSV, TSV}
	case ReferenceBibtexDT:
		return []FileType{BIBTEX}
	case ReferenceJsonDT:
		return []FileType{JSONCSL}
	default:
		return []FileType{UnknownFileType}
	}
}

// ID returns the string ID of the DataType.
func (dt DataType) ID() string {
	return DataTypeToString[dt]
}

// String returns the string representation of the DataType.
func (dt DataType) String() string {
	return ToStr(dt.ID())
}

// StringToDataType maps strings to DataTypes.
var StringToDataType = map[string]DataType{
	"Author":               AuthorDT,
	"NameRelation":         NameRelationDT,
	"Taxon":                TaxonDT,
	"Synonym":              SynonymDT,
	"NameUsage":            NameUsageDT,
	"Name":                 NameDT,
	"TaxonProperty":        TaxonPropertyDT,
	"TaxonConceptRelation": TaxonConceptRelationDT,
	"SpeciesInteraction":   SpeciesInteractionDT,
	"SpeciesEstimate":      SpeciesEstimateDT,
	"Reference":            ReferenceDT,
	"ReferenceJson":        ReferenceJsonDT,
	"ReferenceBibTex":      ReferenceBibtexDT,
	"TypeMaterial":         TypeMaterialDT,
	"Distribution":         DistributionDT,
	"Media":                MediaDT,
	"VernacularName":       VernacularNameDT,
	"Treatment":            TreatmentDT,
}

// DataTypeToString maps DataType to string.
var DataTypeToString = func() map[DataType]string {
	res := make(map[DataType]string)
	for k, v := range StringToDataType {
		res[v] = k
	}
	res[UnkownDT] = "Unknown"
	return res
}()

// LowCaseToDataType provides map of low-case strings to DataType.
var LowCaseToDataType = func() map[string]DataType {
	res := make(map[string]DataType)
	for k, v := range StringToDataType {
		res[strings.ToLower(k)] = v
	}
	return res
}()

// NewDataType creates a new DataType from a string representation.
func NewDataType(s, ext string) DataType {
	s = strings.ToLower(s)
	ext = strings.ToLower(ext)
	// extension can be .json or .jsonl
	if strings.HasPrefix(ext, ".json") {
		s += "Json"
	}
	if dt, ok := LowCaseToDataType[s]; ok {
		return dt
	}
	return UnkownDT
}

type DistrStatus int

const (
	UnknownDistSt DistrStatus = iota
	Native
	Domesticated
	Alien
	Uncertain
)

var distrStatusToString = map[DistrStatus]string{
	Native:       "NATIVE",
	Domesticated: "DOMESTICATED",
	Alien:        "ALIEN",
	Uncertain:    "UNCERTAIN",
}

var stringToDistrStatus = func() map[string]DistrStatus {
	res := make(map[string]DistrStatus)
	for k, v := range distrStatusToString {
		res[v] = k
	}
	return res
}()

// ID return the string ID of DistrStatus.
func (ds DistrStatus) ID() string {
	if res, ok := distrStatusToString[ds]; ok {
		return res
	}
	return ""
}

// String return the string representation of DistrStatus.
func (ds DistrStatus) String() string {
	return ToStr(ds.ID())
}

func NewDistrStatus(s string) DistrStatus {
	s = strings.ToUpper(s)
	if res, ok := stringToDistrStatus[s]; ok {
		return res
	}
	return UnknownDistSt
}

type EstimateType int

const (
	UnknownET EstimateType = iota
	SpeciesExtinct
	SpeciesLiving
	EstimatedSpecies
)

var estTypeToString = map[EstimateType]string{
	SpeciesExtinct:   "SPECIES_EXTINCT",
	SpeciesLiving:    "SPECIES_LIVING",
	EstimatedSpecies: "ESTIMATED_SPECIES",
}

var stringToEstType = func() map[string]EstimateType {
	res := make(map[string]EstimateType)
	for k, v := range estTypeToString {
		res[v] = k
	}
	return res
}()

// ID returns the string ID of EstimatedType.
func (et EstimateType) ID() string {
	if res, ok := estTypeToString[et]; ok {
		return res
	}
	return ""
}

// String returns the string representation of EstimatedType.
func (et EstimateType) String() string {
	return ToStr(et.ID())
}

func NewEstimateType(s string) EstimateType {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "_")
	if res, ok := stringToEstType[s]; ok {
		return res
	}
	return UnknownET
}

type TaxonConceptRelType int

const (
	UnknownTCT TaxonConceptRelType = iota
	Equals
	Includes
	IncludedIn
	Overlaps
	Excludes
)

var tcRelTypeToString = map[TaxonConceptRelType]string{
	Equals:     "EQUALS",
	Includes:   "INCLUDES",
	IncludedIn: "INCLUDED_IN",
	Overlaps:   "OVERLAPS",
	Excludes:   "EXCLUDES",
}

var stringToTCRelType = func() map[string]TaxonConceptRelType {
	res := make(map[string]TaxonConceptRelType)
	for k, v := range tcRelTypeToString {
		res[v] = k
	}
	return res
}()

// ID returns the string ID of TaxonConceptRelType.
func (t TaxonConceptRelType) ID() string {
	if res, ok := tcRelTypeToString[t]; ok {
		return res
	}
	return ""
}

// String returns the string representation of TaxonConceptRelType.
func (t TaxonConceptRelType) String() string {
	return ToStr(t.ID())
}

func NewTaxonConceptRelType(s string) TaxonConceptRelType {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "_")
	if res, ok := stringToTCRelType[s]; ok {
		return res
	}
	return UnknownTCT
}
