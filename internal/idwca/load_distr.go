package idwca

import (
	"context"
	"errors"
	"strings"

	"github.com/gnames/gnuuid"
	"github.com/sfborg/sflib/pkg/coldp"
	dwca "github.com/sfborg/sflib/pkg/dwca"
	"golang.org/x/sync/errgroup"
)

func (a *idwca) LoadDistribution(
	ctx context.Context,
	idx int,
	chOut chan<- coldp.Data,
) error {
	chIn := make(chan []string)
	g, ctx2 := errgroup.WithContext(ctx)

	ext, err := a.extByIdx(idx)
	if err != nil {
		return err
	}

	g.Go(func() error {
		return a.distrWorker(ctx2, ext, chIn, chOut)
	})

	g.Go(func() error {
		defer close(chIn)
		_, err := a.ExtensionStream(ctx2, idx, chIn)
		return err
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *idwca) distrWorker(
	ctx context.Context,
	ext *dwca.Extension,
	chIn <-chan []string,
	chOut chan<- coldp.Data,
) error {
	fieldsMap := fieldsMap(ext.Fields)
	coreID := ext.CoreID.Idx

	data := coldp.Data{
		Distributions: make([]coldp.Distribution, 0, a.cfg.BatchSize),
		References:    make([]coldp.Reference, 0, a.cfg.BatchSize),
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-chIn:
			if !ok {
				if len(data.Distributions) > 0 {
					chOut <- data
				}
				return nil
			}

			if len(data.Distributions) >= a.cfg.BatchSize {
				chOut <- data
				data.Distributions = data.Distributions[:0]
				data.References = data.References[:0]
			}
			distr, ref := a.processDistrRow(v, coreID, fieldsMap)
			if distr == nil {
				continue
			}
			data.Distributions = append(data.Distributions, *distr)
			if ref != nil {
				data.References = append(data.References, *ref)
			}
		}
	}
}

func (a *idwca) processDistrRow(
	row []string,
	coreID int,
	fieldsMap map[string]int,
) (*coldp.Distribution, *coldp.Reference) {
	dist := &coldp.Distribution{}
	var ref *coldp.Reference

	dist.TaxonID = row[coreID]

	// dwc:taxonID
	// dwc:occurrenceStatus
	// dwc:locationID
	// dwc:locality
	// dwc:countryCode
	// dcterms:source
	dist.TaxonID = fieldVal(row, fieldsMap, "taxonid")

	if dist.TaxonID == "" {
		return nil, nil
	}

	status := fieldVal(row, fieldsMap, "occurrencestatus")
	dist.Status = coldp.NewDistrStatus(status)

	dist.Area = fieldVal(row, fieldsMap, "locality")
	country := fieldVal(row, fieldsMap, "countrycode")
	if len(country) == 2 {
		dist.AreaID = country
		dist.Gazetteer = coldp.NewGazetteerEnt("iso")
	} else if country != "" {
		dist.Area = country
	}

	citation := fieldVal(row, fieldsMap, "source")
	if citation != "" {
		dist.ReferenceID = gnuuid.New(citation).String()
		ref = &coldp.Reference{ID: dist.ReferenceID, Citation: citation}
	}

	location := fieldVal(row, fieldsMap, "locationid")
	if location == "" {
		return dist, ref
	}

	data := strings.Split(location, ":")
	if len(data) != 2 {
		return dist, ref
	}

	gazeteer := coldp.NewGazetteerEnt(data[0])
	if gazeteer.ID() == "" {
		return dist, ref
	}

	dist.Gazetteer = gazeteer
	dist.AreaID = data[1]

	return dist, ref
}
