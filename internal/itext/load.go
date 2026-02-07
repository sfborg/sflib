package itext

import (
	"bufio"
	"context"
	"errors"
	"os"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/parser"
	"golang.org/x/sync/errgroup"
)

func (a *itext) Load(
	ctx context.Context,
	ch chan<- coldp.NameUsage,
	jobsNum int,
	nomcode nomcode.Code,
) error {
	a.code = nomcode
	a.jobsNum = jobsNum
	a.parserPool = parser.Pool(jobsNum)

	g, ctx2 := errgroup.WithContext(ctx)
	chIn := make(chan string)

	g.Go(func() error {
		err := a.read(ctx2, chIn, a.sfgaFilePath)
		close(chIn)
		return err
	})

	for range jobsNum {
		g.Go(func() error {
			return a.process(ctx2, chIn, ch)
		})
	}

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (a *itext) read(
	ctx context.Context,
	chIn chan<- string,
	path string,
) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chIn <- scanner.Text():
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

func (a *itext) process(
	ctx context.Context,
	chIn <-chan string,
	chOut chan<- coldp.NameUsage,
) error {
	code := nomcode.Unknown
	switch a.code {
	case nomcode.Botanical, nomcode.Cultivars:
		code = nomcode.Botanical
	}

	p := <-a.parserPool[code]

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-chIn:
			if !ok {
				a.parserPool[code] <- p
				return nil
			}
			chOut <- a.processLine(p, line)
		}
	}
}

func (a *itext) processLine(p gnparser.GNparser, line string) coldp.NameUsage {
	prsd := p.ParseName(line).Flatten()
	res := coldp.NameUsage{
		ID:                        prsd.VerbatimID,
		GlobalID:                  prsd.VerbatimID,
		ScientificName:            line,
		Authorship:                prsd.Authorship,
		ScientificNameString:      line,
		ParseQuality:              coldp.ToInt(prsd.ParseQuality),
		CanonicalSimple:           prsd.CanonicalSimple,
		CanonicalFull:             prsd.CanonicalFull,
		CanonicalStemmed:          prsd.CanonicalStemmed,
		Cardinality:               coldp.ToInt(prsd.Cardinality),
		Virus:                     coldp.ToBool(prsd.Virus),
		Hybrid:                    prsd.Hybrid,
		Surrogate:                 prsd.Surrogate,
		Authors:                   prsd.Authors,
		GnID:                      prsd.VerbatimID,
		Rank:                      coldp.NewRank(prsd.Rank),
		Uninomial:                 prsd.Uninomial,
		GenericName:               prsd.Genus,
		InfragenericEpithet:       prsd.Subgenus,
		SpecificEpithet:           prsd.Species,
		InfraspecificEpithet:      prsd.Infraspecies,
		CultivarEpithet:           prsd.CultivarEpithet,
		CombinationAuthorship:     prsd.CombinationAuthorship,
		CombinationExAuthorship:   prsd.CombinationExAuthorship,
		CombinationAuthorshipYear: prsd.CombinationAuthorshipYear,
		BasionymAuthorship:        prsd.BasionymAuthorship,
		BasionymExAuthorship:      prsd.BasionymExAuthorship,
		BasionymAuthorshipYear:    prsd.BasionymAuthorshipYear,
		Code:                      a.code,
	}

	return res
}
