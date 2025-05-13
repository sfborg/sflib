package idwca

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/gnames/gnuuid"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca"
	"github.com/sfborg/sflib/pkg/dwca/diagn"
	"golang.org/x/sync/errgroup"
)

func (a *idwca) LoadCore(
	ctx context.Context,
	chOut chan<- coldp.Data,
) error {
	chIn := make(chan []string)
	g, ctx2 := errgroup.WithContext(ctx)

	for range a.cfg.JobsNum {
		g.Go(func() error {
			return a.coreWorker(ctx2, chIn, chOut)
		})
	}

	g.Go(func() error {
		defer close(chIn)
		_, err := a.CoreStream(ctx2, chIn)
		return err
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *idwca) coreWorker(
	ctx context.Context,
	chIn chan []string,
	chOut chan<- coldp.Data,
) error {
	fieldsMap := fieldsMap(a.Meta().Core.Fields)
	coreID := a.Meta().Core.ID.Idx

	data := coldp.Data{
		NameUsages: make([]coldp.NameUsage, 0, a.cfg.BatchSize),
		References: make([]coldp.Reference, 0, a.cfg.BatchSize),
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-chIn:
			if !ok {
				if len(data.NameUsages) > 0 {
					chOut <- data
				}
				return nil
			}
			if len(data.NameUsages) >= a.cfg.BatchSize {
				chOut <- data
				data = coldp.Data{
					NameUsages: make([]coldp.NameUsage, 0, a.cfg.BatchSize),
					References: make([]coldp.Reference, 0, a.cfg.BatchSize),
				}
			}

			nu, ref := a.processCoreRow(v, coreID, fieldsMap)
			data.NameUsages = append(data.NameUsages, nu)
			if ref != nil {
				data.References = append(data.References, *ref)
			}
		}
	}
}

func fieldsMap(fields []dwca.Field) map[string]int {
	fieldsMap := make(map[string]int)
	for _, v := range fields {
		term := filepath.Base(v.Term)
		term = strings.ToLower(term)
		_, noPrefix, found := strings.Cut(term, ":")
		if found {
			term = noPrefix
		}
		fieldsMap[term] = v.Idx
	}
	return fieldsMap
}

func fieldVal(
	row []string,
	fielsMap map[string]int,
	name string) string {
	if idx, ok := fielsMap[name]; ok {
		if idx >= len(row) {
			return ""
		}
		return row[idx]
	}
	return ""
}

func (a *idwca) processCoreRow(
	row []string,
	idIdx int,
	fieldsMap map[string]int,
) (coldp.NameUsage, *coldp.Reference) {
	nu := coldp.NameUsage{ID: row[idIdx]}
	nu.SourceID = fieldVal(row, fieldsMap, "datasetid")
	nu.ScientificName = fieldVal(row, fieldsMap, "scientificname")
	nu.Authorship = fieldVal(row, fieldsMap, "scientificnameauthorship")

	nu.ScientificNameString = fieldVal(row, fieldsMap, "scientificnamestring")
	nu.LocalID = fieldVal(row, fieldsMap, "localid")
	nu.GlobalID = fieldVal(row, fieldsMap, "globalid")

	parentID := fieldVal(row, fieldsMap, "parentnameusageid")
	// iNat provides a URL to parentID instead of bare ID :-/
	if strings.HasPrefix(parentID, "http") {
		parentID = filepath.Base(parentID)
	}
	nu.ParentID = parentID

	nomCode := fieldVal(row, fieldsMap, "nomenclaturalcode")
	nu.Code = nomcode.New(nomCode)
	rank := fieldVal(row, fieldsMap, "taxonrank")
	nu.Rank = coldp.NewRank(rank)
	nu = addTaxonomicStatus(nu, row, fieldsMap)

	nu.Notho = coldp.NewNamePart(fieldVal(row, fieldsMap, "notho"))
	nu.Kingdom = fieldVal(row, fieldsMap, "kingdom")
	nu.Phylum = fieldVal(row, fieldsMap, "phylum")
	nu.Class = fieldVal(row, fieldsMap, "class")
	nu.Order = fieldVal(row, fieldsMap, "order")
	nu.Superfamily = fieldVal(row, fieldsMap, "superfamily")
	nu.Family = fieldVal(row, fieldsMap, "family")
	nu.Subfamily = fieldVal(row, fieldsMap, "subfamily")
	nu.Tribe = fieldVal(row, fieldsMap, "tribe")
	nu.Genus = fieldVal(row, fieldsMap, "genericName")
	nu.InfragenericEpithet = fieldVal(row, fieldsMap, "infragenericepither")
	nu.InfraspecificEpithet = fieldVal(row, fieldsMap, "infraspecificepither")
	nu.CultivarEpithet = fieldVal(row, fieldsMap, "cultivarepithet")
	nu.AccordingToID = fieldVal(row, fieldsMap, "nameaccordingto")
	nu.NameStatus = coldp.NewNomStatus(fieldVal(row, fieldsMap, "nomenclaturalstatus"))
	nu.Remarks = fieldVal(row, fieldsMap, "taxonremarks")
	if nu.ScientificNameString == "" {
		a.setNameString(&nu)
	}

	var ref *coldp.Reference
	citation := fieldVal(row, fieldsMap, "namepublishedin")
	if citation != "" {
		nu.ReferenceID = gnuuid.New(citation).String()
		ref = &coldp.Reference{ID: nu.ReferenceID, Citation: citation}
	}

	code := a.getParserCode(nu.Code)
	p := a.parserPool[code].Get().(gnparser.GNparser)
	nu.Amend(p)
	a.parserPool[code].Put(p)

	return nu, ref
}

func (a *idwca) getParserCode(code nomcode.Code) nomcode.Code {
	res := nomcode.Unknown
	botCodes := []nomcode.Code{nomcode.Botanical, nomcode.Cultivars}
	if slices.Contains(botCodes, a.cfg.Code) {
		res = nomcode.Botanical
	}
	if slices.Contains(botCodes, code) {
		res = nomcode.Botanical
	}
	return res
}

func (a *idwca) setNameString(nu *coldp.NameUsage) {
	d := a.Diagnostics()
	switch d.SciNameType {
	case diagn.SciNameFull:
		nu.ScientificNameString = nu.ScientificName
	case diagn.SciNameCanonical:
		name := nu.ScientificName + " " + nu.Authorship
		nu.ScientificNameString = strings.TrimSpace(name)
	default:
	}
}

func addTaxonomicStatus(
	nu coldp.NameUsage,
	row []string,
	fieldMap map[string]int,
) coldp.NameUsage {
	nu.TaxonomicStatus = coldp.UnknownTaxSt

	// if neither acceptedNameUsageID or taxonomicStatus are given
	// we cannot assert taxonomic status
	_, hasStatus := fieldMap["taxonomicstatus"]
	_, hasAccepted := fieldMap["acceptednameusageid"]
	if !hasStatus && !hasAccepted {
		return nu
	}

	// use provided addTaxonomicStatus if possible
	ts := fieldVal(row, fieldMap, "taxonomicstatus")
	nu.TaxonomicStatus = coldp.NewTaxonomicStatus(ts)
	if nu.TaxonomicStatus != coldp.UnknownTaxSt {
		return nu
	}

	acceptedNameUsageID := fieldVal(row, fieldMap, "acceptednameusageid")
	if acceptedNameUsageID != "" && nu.ID != acceptedNameUsageID {
		// for NameUsage's synonyms accepted ID goes to parent, and
		// synonymy is expressed by taxonomic status
		nu.ParentID = acceptedNameUsageID
		nu.TaxonomicStatus = coldp.SynonymTS
		return nu
	}

	nu.TaxonomicStatus = coldp.AcceptedTS
	return nu
}
