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

const coreFile = "Taxon.csv"

func (a *idwca) WriteCore(
	ctx context.Context,
	ch <-chan coldp.NameUsage,
) error {
	path := filepath.Join(a.rootDir, coreFile)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", coreFile, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := dwca.TermHeaders(dwca.CoreTerms)
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("cannot write core header: %w", err)
	}

	var count int
	for nu := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		row := coreRow(nu)
		if err := w.Write(row); err != nil {
			return fmt.Errorf("cannot write core row: %w", err)
		}

		count++
		if count%50_000 == 0 {
			util.Progress(count, "core")
		}
	}
	util.ProgressEnd()

	if count == 0 {
		// Remove the empty file and return an error.
		os.Remove(path)
		return fmt.Errorf("no core records: source contains no name usages")
	}

	a.coreWritten = true
	return nil
}

// coreRow converts a NameUsage into a CSV row matching CoreTerms order.
// In DwCA, synonyms use acceptedNameUsageID (not parentNameUsageID),
// while accepted taxa use parentNameUsageID.
func coreRow(nu coldp.NameUsage) []string {
	var parentID, acceptedID string
	if isSynonym(nu.TaxonomicStatus) {
		acceptedID = nu.ParentID
	} else {
		parentID = nu.ParentID
	}

	return []string{
		nu.ID,
		parentID,
		acceptedID,
		nu.BasionymID,
		nu.TaxonomicStatus.String(),
		nu.Rank.String(),
		nu.ScientificName,
		nu.Authorship,
		nu.GenericName,
		nu.InfragenericEpithet,
		nu.SpecificEpithet,
		nu.InfraspecificEpithet,
		nu.CultivarEpithet,
		nu.Code.String(),
		nu.NameStatus.String(),
		nu.Realm,
		nu.Kingdom,
		nu.Phylum,
		nu.Subphylum,
		nu.Class,
		nu.Subclass,
		nu.Order,
		nu.Suborder,
		nu.Superfamily,
		nu.Family,
		nu.Subfamily,
		nu.Tribe,
		nu.Subtribe,
		nu.Genus,
		nu.Subgenus,
		nu.Section,
		nu.Species,
		nu.Remarks,
		nu.Link,
		nu.Modified,
	}
}

func isSynonym(ts coldp.TaxonomicStatus) bool {
	switch ts {
	case coldp.SynonymTS, coldp.AmbiguousSynonymTS, coldp.MisappliedTS:
		return true
	}
	return false
}
