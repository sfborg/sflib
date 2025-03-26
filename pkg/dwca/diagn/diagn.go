package diagn

import (
	"log/slog"
	"strings"

	"github.com/gnames/gnparser"
)

type Diagnostics struct {
	SciNameType
	SynonymType
	HierType
}

func New(
	d []map[string]string,
	exts map[string]string,
) *Diagnostics {
	res := Diagnostics{}
	slog.Info("Checking if ScientificName contains authorship.")
	res.SciNameType = sciNameType(d)
	slog.Info("Finding where synonyms are.")
	res.SynonymType = synonymType(d, exts)
	slog.Info("Checking if hierarchy is flat, parent/child or none.")
	res.HierType = hierType(d)
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
		if k == "kingdom" || k == "phylum" ||
			k == "class" || k == "order" || k == "family" ||
			k == "genus" || k == "species" {
			count++
		}
		if k == "parentnameusageid" || k == "highertaxonid" {
			isTree = true
		}
	}
	if count > 5 {
		return HierFlat
	}

	if isTree {
		return HierTree
	}

	return HierUnknown
}

func synonymType(
	cs []map[string]string,
	exts map[string]string,
) SynonymType {
	if len(cs) == 0 {
		return SynUnknown
	}

	for k, v := range exts {
		if strings.HasPrefix(k, "synonym") || strings.HasPrefix(v, "synonym") {
			return SynExtension
		}
	}

	if _, ok := cs[0]["acceptednameusageid"]; ok {
		return SynAcceptedID
	}

	for _, v := range cs {
		isHier := checkHierarchy(v)
		if isHier {
			return SynHierarchy
		}
	}

	return SynUnknown
}

func checkHierarchy(v map[string]string) bool {
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
		return SciNameUnknown
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
			if compositeNum > 100 {
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
			if canonicalNum > 100 {
				break
			}
			continue
		}

		if parsed.Authorship != nil {
			fullNum++
			if fullNum > 100 {
				break
			}
			continue
		}
	}

	if fullNum > 0 && canonicalNum+compositeNum == 0 {
		return SciNameFull
	}
	if canonicalNum > 0 && fullNum+compositeNum == 0 {
		return SciNameCanonical
	}
	if compositeNum > 0 && canonicalNum+fullNum == 0 {
		return SciNameComposite
	}
	return SciNameUnknown
}
