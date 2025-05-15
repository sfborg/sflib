package ixsv

import (
	"context"
	"errors"
	"strings"

	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnfmt/gncsv/config"
	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/parser"
	"golang.org/x/sync/errgroup"
)

func (a *ixsv) Load(
	ctx context.Context,
	ch chan<- coldp.NameUsage,
	jobsNum int,
	nomCode nomcode.Code,
) error {
	opt := config.OptPath(a.filePath)
	cfg, err := config.New(opt)
	if err != nil {
		return err
	}
	a.reader = gncsv.New(cfg)
	a.headers = coldp.NormalizeHeaders(a.reader.Headers())
	a.code = nomCode
	a.jobsNum = jobsNum
	a.parserPool = parser.Pool(jobsNum)

	g, ctx2 := errgroup.WithContext(ctx)
	chIn := make(chan []string)

	g.Go(func() error {
		defer close(chIn)
		_, err := a.reader.Read(ctx2, chIn)
		return err
	})

	for range jobsNum {
		g.Go(func() error {
			return a.process(ctx2, chIn, ch)
		})
	}

	err = g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *ixsv) process(
	ctx context.Context,
	chIn <-chan []string,
	chOut chan<- coldp.NameUsage,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case row, ok := <-chIn:
			if !ok {
				return nil
			}
			a.processRow(row, chOut)
		}
	}
}

func (a *ixsv) processRow(
	row []string,
	chOut chan<- coldp.NameUsage,
) {
	rowCode := nomcode.New(a.getVal(row, "code"))
	code := parser.ParserCode(a.cfg.Code, rowCode)

	p := a.parserPool[code].Get().(gnparser.GNparser)

	nu := a.getNameUsage(p, row)
	a.parserPool[code].Put(p)

	chOut <- nu
}

func (a *ixsv) getVal(row []string, field string) string {
	if idx, ok := a.headers[field]; ok {
		res := row[idx]
		return strings.TrimSpace(res)
	}
	return ""
}

func (a *ixsv) getNameUsage(
	p gnparser.GNparser,
	row []string,
) coldp.NameUsage {
	res := coldp.NameUsage{
		ID:                a.getVal(row, "id"),
		AlternativeID:     a.getVal(row, "alternativeid"),
		NameAlternativeID: a.getVal(row, "namealternativeid"),
		LocalID:           a.getVal(row, "localid"),
		GlobalID:          a.getVal(row, "globalid"),
		SourceID:          a.getVal(row, "sourceid"),
		ParentID:          a.getVal(row, "parentid"),
		BasionymID:        a.getVal(row, "basyonymid"),
		TaxonomicStatus: coldp.NewTaxonomicStatus(
			a.getVal(row, "taxonomicstatus"),
		),
		ScientificName:            a.getVal(row, "scientificname"),
		Authorship:                a.getVal(row, "authorship"),
		ScientificNameString:      a.getVal(row, "scientificnamestring"),
		Rank:                      coldp.NewRank(a.getVal(row, "rank")),
		Notho:                     coldp.NewNamePart(a.getVal(row, "notho")),
		Uninomial:                 a.getVal(row, "uninomial"),
		GenericName:               a.getVal(row, "genericname"),
		InfragenericEpithet:       a.getVal(row, "infragenericepithet"),
		SpecificEpithet:           a.getVal(row, "specificepithet"),
		InfraspecificEpithet:      a.getVal(row, "infraspecificepithet"),
		CultivarEpithet:           a.getVal(row, "cultivarepithet"),
		CombinationAuthorship:     a.getVal(row, "combinationauthorship"),
		CombinationAuthorshipID:   a.getVal(row, "combinationauthorshipid"),
		CombinationExAuthorship:   a.getVal(row, "combinationexauthorship"),
		CombinationExAuthorshipID: a.getVal(row, "combinationexauthorshipid"),
		CombinationAuthorshipYear: a.getVal(row, "combinationauthorshipyear"),
		BasionymAuthorship:        a.getVal(row, "basionymauthorship"),
		BasionymAuthorshipID:      a.getVal(row, "basionymauthorshipid"),
		BasionymExAuthorship:      a.getVal(row, "basionymexauthorship"),
		BasionymExAuthorshipID:    a.getVal(row, "basionymexauthorshipid"),
		BasionymAuthorshipYear:    a.getVal(row, "basionymauthorshipyear"),
		NamePhrase:                a.getVal(row, "namephrase"),
		NameReferenceID:           a.getVal(row, "namereferenceid"),
		PublishedInYear:           a.getVal(row, "publishedinyear"),
		PublishedInPage:           a.getVal(row, "publishedinpage"),
		PublishedInPageLink:       a.getVal(row, "publishedinpagelink"),
		Gender:                    coldp.NewGender(a.getVal(row, "gender")),
		Etymology:                 a.getVal(row, "etymology"),
		Code:                      nomcode.New(a.getVal(row, "code")),
		NameStatus: coldp.NewNomStatus(
			a.getVal(row, "namestatus"),
		),
		AccordingToID:       a.getVal(row, "accordingtoid"),
		AccordingToPage:     a.getVal(row, "accordingtopage"),
		AccordingToPageLink: a.getVal(row, "accordingtopagelink"),
		ReferenceID:         a.getVal(row, "referenceid"),
		Scrutinizer:         a.getVal(row, "scrutinizer"),
		ScrutinizerID:       a.getVal(row, "scrutinizerid"),
		ScrutinizerDate:     a.getVal(row, "scrutinizerdate"),
		Species:             a.getVal(row, "species"),
		Section:             a.getVal(row, "section"),
		Subgenus:            a.getVal(row, "subgenus"),
		Genus:               a.getVal(row, "genus"),
		Subtribe:            a.getVal(row, "subtribe"),
		Tribe:               a.getVal(row, "tribe"),
		Subfamily:           a.getVal(row, "subfamily"),
		Family:              a.getVal(row, "family"),
		Superfamily:         a.getVal(row, "superfamily"),
		Suborder:            a.getVal(row, "suborder"),
		Order:               a.getVal(row, "order"),
		Subclass:            a.getVal(row, "subclass"),
		Class:               a.getVal(row, "class"),
		Subphylum:           a.getVal(row, "subphylum"),
		Phylum:              a.getVal(row, "phylum"),
		Kingdom:             a.getVal(row, "kingdom"),
		Link:                a.getVal(row, "link"),
		NameRemarks:         a.getVal(row, "nameremarks"),
		Remarks:             a.getVal(row, "remarks"),
		Modified:            a.getVal(row, "modified"),
		ModifiedBy:          a.getVal(row, "modifiedby"),
	}

	// add parsed data
	res.ScientificNameString = res.ScientificName

	if res.Authorship != "" &&
		!strings.HasSuffix(res.ScientificName, res.Authorship) {
		res.ScientificNameString += " " + res.Authorship
	}
	res.Amend(p)
	return res
}
