package coldp

import "strings"

// Rank represents the taxonomic rank of a scientific name.
type Rank int

// ID returns the ID of the rank.
func (r Rank) ID() string {
	if res, ok := rankToString[r]; ok {
		return res
	}
	return ""
}

// String returns the string representation of the rank.
func (r Rank) String() string {
	return ToStr(r.ID())
}

// NewRank creates a new Rank from a string representation.
// It handles trimming and normalization to ensure consistent matching.
func NewRank(s string) Rank {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".")
	s = strings.ToUpper(s)
	if s == "" {
		return Unranked
	}

	var res Rank
	var ok bool
	if res, ok = abbrToRank[s]; ok {
		return res
	}
	if res, ok = stringToRank[s]; ok {
		return res
	}
	return UnknownRank
}

// Constants for different taxonomic ranks.
const (
	UnknownRank Rank = iota
	Superdomain
	Domain
	Subdomain
	Infradomain
	Empire
	Realm
	Subrealm
	Superkingdom
	Kingdom
	Subkingdom
	Infrakingdom
	Superphylum
	Phylum
	Subphylum
	Infraphylum
	Parvphylum
	Microphylum
	Nanophylum
	Claudius
	Gigaclass
	Megaclass
	Superclass
	Class
	subclass
	Infraclass
	Subterclass
	Parvclass
	Superdivision
	Division
	Subdivision
	Infradivision
	Superlegion
	Legion
	Sublegion
	Infralegion
	Megacohort
	Supercohort
	Cohort
	Subcohort
	Infracohort
	Gigaorder
	Magnorder
	Grandorder
	Mirorder
	Superorder
	Order
	Nanorder
	Hypoorder
	Minorder
	Suborder
	Infraorder
	Parvorder
	SupersectionZoology
	SectionZoology
	SubsectionZoology
	Falanx
	Gigafamily
	Megafamily
	Grandfamily
	Superfamily
	Epifamily
	Family
	Subfamily
	Infrafamily
	Supertribe
	Tribe
	Subtribe
	Infratribe
	SupragenericName
	Supergenus
	Genus
	Subgenus
	Infragenus
	SupersectionBotany
	SectionBotany
	SubsectionBotany
	Superseries
	Series
	Subseries
	InfragenericName
	SpeciesAggregate
	Species
	InfraspecificName
	Grex
	Klepton
	Subspecies
	CultivarGroup
	Convariety
	InfrasubspecificName
	Proles
	Natio
	Aberration
	Morph
	Supervariety
	Variety
	Subvariety
	Superform
	Form
	Subform
	Pathovar
	Biovar
	Chemovar
	Morphovar
	Phagovar
	Serovar
	Chemoform
	FormaSpecialis
	Lusus
	Cultivar
	Mutatio
	Strain
	Other
	Unranked
	Subclass
	Section
	Supersection
	Subsection
)

// abbrToRank maps rank abbreviations to their corresponding Rank values.
var abbrToRank = map[string]Rank{
	"SUPERDOM":   Superdomain,
	"DOM":        Domain,
	"SUPERREG":   Superkingdom,
	"REG":        Kingdom,
	"SUBREG":     Subkingdom,
	"INFRAREG":   Infrakingdom,
	"SUPERPHYL":  Superphylum,
	"PHYL":       Phylum,
	"SUBPHYL":    Subphylum,
	"INFRAPHYL":  Infraphylum,
	"PARVPHYL":   Parvphylum,
	"MICROPHYL":  Microphylum,
	"GIGACL":     Gigaclass,
	"MEGACL":     Megaclass,
	"SUPERCL":    Superclass,
	"CL":         Class,
	"SUBCL":      Subclass,
	"SUBTERCL":   Subterclass,
	"PARVCL":     Parvclass,
	"SUPERDIV":   Superdivision,
	"DIV":        Division,
	"SUBDIV":     Subdivision,
	"INFRADIV":   Infradivision,
	"SUPERLEG":   Superlegion,
	"LEGION":     Legion,
	"SUBLEG":     Sublegion,
	"INFRALEG":   Infralegion,
	"GIGAORD":    Gigaorder,
	"GRANDORD":   Grandorder,
	"MIRORD":     Mirorder,
	"SUPERORD":   Superorder,
	"ORD":        Order,
	"NANORD":     Nanorder,
	"HYPORD":     Hypoorder,
	"MINORD":     Minorder,
	"SUBORD":     Suborder,
	"INFRAORD":   Infraorder,
	"PARVORD":    Parvorder,
	"SUPERSECT":  Supersection,
	"SECT":       Section,
	"SUBSECT":    Subsection,
	"MEGAFAM":    Megafamily,
	"GRANDFAM":   Grandfamily,
	"SUPERFAM":   Superfamily,
	"FAM":        Family,
	"SUBFAM":     Subfamily,
	"INFRAFAM":   Infrafamily,
	"SUPERTRIB":  Supertribe,
	"TRIB":       Tribe,
	"SUBTRIB":    Subtribe,
	"INFRATRIB":  Infratribe,
	"SUPERGEN":   Supergenus,
	"GEN":        Genus,
	"SUBGEN":     Subgenus,
	"INFRAGEN":   Infragenus,
	"SUPERSER":   Superseries,
	"SER":        Series,
	"SUBSER":     Subseries,
	"SP":         Species,
	"INFRASP":    InfraspecificName,
	"GX":         Grex,
	"SUBSP":      Subspecies,
	"CONVAR":     Convariety,
	"INFRASUBSP": InfrasubspecificName,
	"AB":         Aberration,
	"SUPERVAR":   Supervariety,
	"VAR":        Variety,
	"SUBVAR":     Subvariety,
	"SUPERF":     Superform,
	"F":          Form,
	"SUBF":       Subform,
	"PV":         Pathovar,
	"F.SP":       FormaSpecialis,
	"CV":         Cultivar,
	"MUT":        Mutatio,
}

// rankToAbbr maps Rank values to their corresponding abbreviations.
var rankToAbbr = func() map[Rank]string {
	res := make(map[Rank]string)
	for k, v := range abbrToRank {
		res[v] = k
	}
	return res
}()

// rankToString maps Rank values to their full string names.
var rankToString = map[Rank]string{
	Unranked:             "UNRANKED",
	Aberration:           "ABERRATION",
	Biovar:               "BIOVAR",
	Chemoform:            "CHEMOFORM",
	Chemovar:             "CHEMOVAR",
	Class:                "CLASS",
	Cohort:               "COHORT",
	Convariety:           "CONVARIETY",
	Cultivar:             "CULTIVAR",
	CultivarGroup:        "CULTIVAR_GROUP",
	Division:             "DIVISION",
	Domain:               "DOMAIN",
	Epifamily:            "EPIFAMILY",
	Falanx:               "FALANX",
	Family:               "FAMILY",
	Form:                 "FORM",
	FormaSpecialis:       "FORMA_SPECIALIS",
	Genus:                "GENUS",
	Gigaclass:            "GIGACLASS",
	Gigaorder:            "GIGAORDER",
	Grandfamily:          "GRANDFAMILY",
	Grandorder:           "GRANDORDER",
	Grex:                 "GREX",
	Hypoorder:            "HYPOORDER",
	Infraclass:           "INFRACLASS",
	Infracohort:          "INFRACOHORT",
	Infradivision:        "INFRADIVISION",
	Infrafamily:          "INFRAFAMILY",
	InfragenericName:     "INFRAGENERIC_NAME",
	Infragenus:           "INFRAGENUS",
	Infrakingdom:         "INFRAKINGDOM",
	Infralegion:          "INFRALEGION",
	Infraorder:           "INFRAORDER",
	Infraphylum:          "INFRAPHYLUM",
	InfraspecificName:    "INFRASPECIFIC_NAME",
	InfrasubspecificName: "INFRASUBSPECIFIC_NAME",
	Infratribe:           "INFRATRIBE",
	Kingdom:              "KINGDOM",
	Klepton:              "KLEPTON",
	Legion:               "LEGION",
	Lusus:                "LUSUS",
	Magnorder:            "MAGNORDER",
	Megaclass:            "MEGACLASS",
	Megacohort:           "MEGACOHORT",
	Megafamily:           "MEGAFAMILY",
	Microphylum:          "MICROPHYLUM",
	Minorder:             "MINORDER",
	Mirorder:             "MIRORDER",
	Morph:                "MORPH",
	Morphovar:            "MORPHOVAR",
	Mutatio:              "MUTATIO",
	Nanophylum:           "NANOPHYLUM",
	Nanorder:             "NANORDER",
	Natio:                "NATIO",
	Order:                "ORDER",
	Other:                "OTHER",
	Parvclass:            "PARVCLASS",
	Parvorder:            "PARVORDER",
	Parvphylum:           "PARVPHYLUM",
	Pathovar:             "PATHOVAR",
	Phagovar:             "PHAGOVAR",
	Phylum:               "PHYLUM",
	Proles:               "PROLES",
	Realm:                "REALM",
	Section:              "SECTION",
	Series:               "SERIES",
	Serovar:              "SEROVAR",
	Species:              "SPECIES",
	SpeciesAggregate:     "SPECIES_AGGREGATE",
	Strain:               "STRAIN",
	Subclass:             "SUBCLASS",
	Subcohort:            "SUBCOHORT",
	Subdivision:          "SUBDIVISION",
	Subfamily:            "SUBFAMILY",
	Subform:              "SUBFORM",
	Subgenus:             "SUBGENUS",
	Subkingdom:           "SUBKINGDOM",
	Sublegion:            "SUBLEGION",
	Suborder:             "SUBORDER",
	Subphylum:            "SUBPHYLUM",
	Subrealm:             "SUBREALM",
	Subsection:           "SUBSECTION",
	Subseries:            "SUBSERIES",
	Subspecies:           "SUBSPECIES",
	Subterclass:          "SUBTERCLASS",
	Subtribe:             "SUBTRIBE",
	Subvariety:           "SUBVARIETY",
	Superclass:           "SUPERCLASS",
	Supercohort:          "SUPERCOHORT",
	Superdivision:        "SUPERDIVISION",
	Superdomain:          "SUPERDOMAIN",
	Superfamily:          "SUPERFAMILY",
	Superform:            "SUPERFORM",
	Supergenus:           "SUPERGENUS",
	Superkingdom:         "SUPERKINGDOM",
	Superlegion:          "SUPERLEGION",
	Superorder:           "SUPERORDER",
	Superphylum:          "SUPERPHYLUM",
	Supersection:         "SUPERSECTION",
	Superseries:          "SUPERSERIES",
	Supertribe:           "SUPERTRIBE",
	Supervariety:         "SUPERVARIETY",
	SupragenericName:     "SUPRAGENERIC_NAME",
	Tribe:                "TRIBE",
	Variety:              "VARIETY",
}

// abbrToRank maps rank strings to their corresponding Rank values.
var stringToRank = func() map[string]Rank {
	res := make(map[string]Rank)
	for k, v := range rankToString {
		res[v] = k
	}
	return res
}()
