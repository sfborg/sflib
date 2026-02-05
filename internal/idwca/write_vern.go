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

const vernFile = "VernacularName.csv"

func (a *idwca) WriteVernaculars(
	ctx context.Context,
	ch <-chan coldp.Vernacular,
) error {
	path := filepath.Join(a.rootDir, vernFile)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", vernFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := dwca.TermHeaders(dwca.VernacularTerms)
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("cannot write vernacular header: %w", err)
	}

	var count int
	for v := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row := vernRow(v)
		if err := w.Write(row); err != nil {
			return fmt.Errorf("cannot write vernacular row: %w", err)
		}

		count++
		if count%50_000 == 0 {
			util.Progress(count, "vernacular")
		}
	}
	util.ProgressEnd()

	if count == 0 {
		// No data — remove the empty file, skip extension.
		os.Remove(path)
		return nil
	}

	a.vernWritten = true
	return nil
}

// vernRow converts a Vernacular into a CSV row matching VernacularTerms order.
func vernRow(v coldp.Vernacular) []string {
	return []string{
		v.TaxonID,
		v.Name,
		v.Language,
		v.Country,
		v.ReferenceID,
	}
}
