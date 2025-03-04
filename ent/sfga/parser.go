package sfga

import (
	"strings"

	"github.com/gnames/gnparser/ent/parsed"
)

// FlatParsed represents the parsed scientific name data.
type FlatParsed struct {
	// Quality indicates the quality score of the parsed name.
	Quality int
	// NameID is a UUID v5 generated from the name verbatim string.
	NameID string
	// CanonicalFull is the full canonical form of the name.
	CanonicalFull string
	// CanonicalSimple is the simple canonical form of the name.
	// It does not have ranks, hybrid signs etc.
	CanonicalSimple string
	// CanonicalStemmed are stemmed CanonicalSimple names where
	// suffixes of specific and infraspecific epithets are removed.
	CanonicalStemmed string
	// Authorship is the authorship of the name.
	Authorship string
	// CombinationAuthorship is the authorship of the combination.
	CombinationAuthorship string
	// Uninomial is the uninomial part of the name.
	Uninomial string
	// Genus is the genus part of the name.
	Genus string
	// Subgenus is the subgenus part of the name.
	Subgenus string
	// Species is the species part of the name.
	Species string
	// Rank of the name.
	Rank string
	// Infraspecies is the infraspecies part of the name.
	Infraspecies string
	// UnparsedTail is the unparsed tail of the name.
	UnparsedTail string
}

func ToFlatParsed(parsedName parsed.Parsed) FlatParsed {
	result := FlatParsed{NameID: parsedName.VerbatimID}
	if !parsedName.Parsed {
		return result
	}

	result = FlatParsed{
		NameID:          parsedName.VerbatimID,
		Quality:         parsedName.ParseQuality,
		CanonicalFull:   parsedName.Canonical.Full,
		CanonicalSimple: parsedName.Canonical.Simple,
	}

	if parsedName.Authorship != nil {
		result.Authorship = parsedName.Authorship.Verbatim
	}

	if parsedName.Authorship != nil && parsedName.Authorship.Combination != nil {
		result.CombinationAuthorship = formatAuthors(parsedName.Authorship.Combination.Authors)
	}

	switch detail := parsedName.Details.(type) {
	case parsed.DetailsUninomial:
		result.Uninomial = detail.Uninomial.Value
	case parsed.DetailsSpecies:
		result.Genus = detail.Species.Genus
		result.Species = detail.Species.Species
	case parsed.DetailsInfraspecies:
		if len(detail.Infraspecies.Infraspecies) == 1 {
			result.Genus = detail.Infraspecies.Genus
			result.Species = detail.Infraspecies.Species.Species
			result.Rank = detail.Infraspecies.Infraspecies[0].Rank
			result.Infraspecies = detail.Infraspecies.Infraspecies[0].Value
		}
	}
	return result
}

func formatAuthors(authorship []string) string {
	result := ""
	switch len(authorship) {

	case 1:
		result = authorship[0]
	case 2:
		result = strings.Join(authorship, " & ")
	default:
		result = strings.Join(authorship[0:len(authorship)-1], ", ")
		result = result + " & " + authorship[len(authorship)-1]
	}
	return result
}
