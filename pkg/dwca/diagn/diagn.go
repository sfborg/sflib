package diagn

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/gnames/gnparser"
)

type Diagnostics struct {
	SciNameType
	SynonymType
	HierType
}

func New(
	// data from core file
	data []map[string]string,
	exts map[string]string,
) *Diagnostics {
	res := Diagnostics{}
	slog.Info("Checking if ScientificName contains authorship.")
	res.SciNameType = sciNameType(data)
	slog.Info("Finding where synonyms are.")
	res.SynonymType = synonymType(data, exts)
	slog.Info("Checking if hierarchy is flat, parent/child, both or none.")
	res.HierType = hierType(data)
	return &res
}

func hierType(d []map[string]string) HierType {
	if len(d) == 0 {
		return HierUnknown
	}

	row := d[0]
	var count int
	var isTree bool
	for k := range row {
		if isRank(k) {
			count++
		}
		if k == "parentnameusageid" || k == "highertaxonid" {
			isTree = true
		}
	}

	if isTree && count > 3 {
		return HierBoth
	}
	if isTree {
		return HierTree
	}
	if count > 3 {
		return HierFlat
	}

	return HierUnknown
}

func isRank(fld string) bool {
	ranks := []string{
		"kingdom", "phylum", "class", "order", "superfamily",
		"family", "subfamily", "genus", "subgenus",
		"specificepithet", "infraspecificepithet",
	}
	if slices.Contains(ranks, fld) {
		return true
	}
	return false
}

func synonymType(
	coreData []map[string]string,
	exts map[string]string,
) SynonymType {
	if len(coreData) == 0 {
		return SynUnknown
	}

	for k, v := range exts {
		if strings.HasPrefix(k, "synonym") || strings.HasPrefix(v, "synonym") {
			return SynExtension
		}
	}

	if _, ok := coreData[0]["acceptednameusageid"]; ok {
		return SynAcceptedID
	}

	for _, v := range coreData {
		synHierarchy := checkSynHierarchy(v)
		if synHierarchy {
			return SynHierarchy
		}
	}

	return SynNone
}

func checkSynHierarchy(v map[string]string) bool {
	st := v["taxonomicstatus"]
	var syn bool
	for _, k := range []string{"synonym", "miss", "invalid", "unavailable"} {
		if strings.Contains(st, k) {
			syn = true
			break
		}
	}
	if !syn {
		return false
	}

	if v["parentnameusageid"] != "" || v["highertaxonid"] != "" {
		return true
	}

	return false
}

func sciNameType(d []map[string]string) SciNameType {
	p := gnparser.New(gnparser.NewConfig())

	if len(d) == 0 {
		slog.Warn("Cannot find scientific name data")
		return SciNameUnknown
	}
	count := 100
	if len(d) < count {
		// we decrease the number assuming that at at least half
		// of existing records should give a hint of the name type.
		count = len(d) / 2
	}

	row0 := d[0]
	if _, ok := row0["scientificname"]; ok {
		if _, ok := row0["scientificnameauthorship"]; !ok {
			return SciNameFull
		}
	}

	var canonicalNum, fullNum, compositeNum int
	for _, v := range d {
		if v["scientificname"] == "" && v["specificepithet"] != "" {
			compositeNum++
			if compositeNum > count {
				break
			}
			continue
		}
		parsed := p.ParseName(v["scientificname"])
		if !parsed.Parsed {
			continue
		}

		authField := strings.TrimSpace(v["scientificnameauthorship"])
		if parsed.Authorship == nil && authField != "" {
			canonicalNum++
			continue
		}

		if parsed.Authorship != nil {
			fullNum++
			continue
		}
	}

	if compositeNum > max(canonicalNum, fullNum) {
		return SciNameComposite
	}

	if fullNum > max(canonicalNum, compositeNum) {
		return SciNameFull
	}
	if canonicalNum > max(fullNum, compositeNum) {
		return SciNameCanonical
	}
	slog.Warn("Cannot determine scientific name type")
	return SciNameUnknown
}

func max(a, b int) int {
	res := a
	if b > a {
		res = b
	}
	return res
}
