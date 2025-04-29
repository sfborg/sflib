package idwca

import (
	"context"
	"errors"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca"
	"golang.org/x/sync/errgroup"
)

// LoadVernacular loads vernacular names from the given extension.
func (a *idwca) LoadVernacular(
	ctx context.Context,
	idx int,
	ext *dwca.Extension,
	chOut chan<- []coldp.Vernacular,
) error {
	chIn := make(chan []string)
	g, ctx2 := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.vernWorker(ctx2, ext, chIn, chOut)
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

func (a *idwca) vernWorker(
	ctx context.Context,
	ext *dwca.Extension,
	chIn <-chan []string,
	chOut chan<- []coldp.Vernacular,
) error {
	fieldsMap := fieldsMap(ext.Fields)
	coreID := ext.CoreID.Idx

	batch := make([]coldp.Vernacular, 0, a.cfg.BatchSize)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-chIn:
			if !ok {
				if len(batch) > 0 {
					chOut <- batch
				}
				return nil
			}

			vrn := a.processVernRow(v, coreID, fieldsMap)
			if vrn == nil {
				continue
			}
			if len(batch) == a.cfg.BatchSize {
				chOut <- batch
				batch = batch[:0]
			}
			batch = append(batch, *vrn)
		}
	}
}

func (a *idwca) processVernRow(
	row []string,
	coreID int,
	fieldsMap map[string]int,
) *coldp.Vernacular {
	var res coldp.Vernacular
	res.TaxonID = row[coreID]

	res.Name = fieldVal(row, fieldsMap, "vernacularname")

	if res.Name == "" {
		return nil
	}

	lang := fieldVal(row, fieldsMap, "language")
	if len(lang) > 3 {
		res.Language = gnlib.LangCode(lang)
	} else {
		res.Language = lang
	}
	res.Area = fieldVal(row, fieldsMap, "locality")
	res.Country = fieldVal(row, fieldsMap, "countrycode")
	return &res
}
