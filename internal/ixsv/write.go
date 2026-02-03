package ixsv

import (
	"context"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *ixsv) Write(
	ctx context.Context,
	ch <-chan coldp.NameUsage,
	filePath string,
) error {
	a.filePath = filePath

	w, err := util.NewWriter(filePath, ',')
	if err != nil {
		return err
	}
	defer w.Close()

	// Write headers
	var nu coldp.NameUsage
	if err := w.Write(nu.CoreHeaders()); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case nu, ok := <-ch:
			if !ok {
				return nil
			}
			// Skip synonyms and misapplied names
			if nu.TaxonomicStatus == coldp.SynonymTS ||
				nu.TaxonomicStatus == coldp.AmbiguousSynonymTS ||
				nu.TaxonomicStatus == coldp.MisappliedTS {
				continue
			}
			if err := w.Write(nu.CoreRow()); err != nil {
				return err
			}
			w.Count++
		}
	}
}
