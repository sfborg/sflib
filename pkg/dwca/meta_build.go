package dwca

import (
	"strconv"
)

// BuildMeta constructs a Meta struct describing the archive contents.
// The core (Taxon) is always included. Vernacular and Distribution
// extensions are included only when the corresponding flag is true.
func BuildMeta(hasVern, hasDistr bool) *Meta {
	m := &Meta{
		EMLFile: "eml.xml",
		Core:    buildCore(),
	}

	if hasVern {
		m.Extensions = append(m.Extensions, buildExtension(
			VernacularNameRowType,
			"VernacularName.csv",
			VernacularTerms,
		))
	}
	if hasDistr {
		m.Extensions = append(m.Extensions, buildExtension(
			DistributionRowType,
			"Distribution.csv",
			DistributionTerms,
		))
	}

	return m
}

func buildCore() *Core {
	return &Core{
		ID:   ID{Index: "0"},
		Attr: csvAttr(TaxonRowType, "Taxon.csv", CoreTerms),
	}
}

func buildExtension(rowType, filename string, terms []DwcTerm) *Extension {
	return &Extension{
		CoreID: CoreID{Index: "0"},
		Attr:   csvAttr(rowType, filename, terms),
	}
}

func csvAttr(rowType, filename string, terms []DwcTerm) *Attr {
	fields := make([]Field, len(terms))
	for i, t := range terms {
		fields[i] = Field{
			Index: strconv.Itoa(i),
			Term:  t.URI,
		}
	}
	return &Attr{
		Encoding:           "utf-8",
		FieldsTerminatedBy: ",",
		LinesTerminatedBy:  "\\n",
		FieldsEnclosedBy:   `"`,
		IgnoreHeaderLines:  "1",
		RowType:            rowType,
		Files:              Files{Locations: []string{filename}},
		Fields:             fields,
	}
}
