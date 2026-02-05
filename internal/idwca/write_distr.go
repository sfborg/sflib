package idwca

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca"
)

const distrFile = "Distribution.csv"

func (a *idwca) WriteDistributions(
	ctx context.Context,
	ch <-chan coldp.Distribution,
) error {
	path := filepath.Join(a.rootDir, distrFile)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", distrFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := dwca.TermHeaders(dwca.DistributionTerms)
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("cannot write distribution header: %w", err)
	}

	var count int
	for d := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row := distrRow(d)
		if err := w.Write(row); err != nil {
			return fmt.Errorf("cannot write distribution row: %w", err)
		}

		count++
		if count%50_000 == 0 {
			util.Progress(count, "distribution")
		}
	}
	util.ProgressEnd()

	if count == 0 {
		// No data — remove the empty file, skip extension.
		os.Remove(path)
		return nil
	}

	a.distrWritten = true
	return nil
}

// distrRow converts a Distribution into a CSV row matching
// DistributionTerms order.
func distrRow(d coldp.Distribution) []string {
	return []string{
		d.TaxonID,
		d.Status.String(),
		d.AreaID,
		d.Area,
		"", // countryCode — not in coldp.Distribution directly
		d.ReferenceID,
	}
}
