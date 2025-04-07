package isfga

import (
	"context"
	"fmt"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadTaxa(
	ctx context.Context,
	ch chan<- coldp.Taxon,
) error {
	q := `
SELECT
	col__id, col__alternative_id, col__source_id, col__parent_id, col__ordinal,
	col__branch_length, col__name_id, col__name_phrase, col__according_to_id,
	col__according_to_page, col__according_to_page_link, col__scrutinizer,
	col__scrutinizer_id, col__status_id, col__extinct,
	col__temporal_range_start_id, col__temporal_range_end_id,
	col__environment_id, col__species, col__section, col__subgenus, col__genus,
	col__subtribe, col__tribe, col__subfamily, col__family, col__superfamily,
	col__suborder, col__order, col__subclass, col__class, col__subphylum,
	col__phylum, col__kingdom, col__reference_id, col__link, col__remarks,
	col__modified, col__modified_by
FROM taxon
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA taxon: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var tx coldp.Taxon

		var status, start, end, env string
		err = rows.Scan(
			&tx.ID, &tx.AlternativeID, &tx.SourceID, &tx.ParentID, &tx.Ordinal,
			&tx.BranchLength, &tx.NameID, &tx.NamePhrase, &tx.AccordingToID,
			&tx.AccordingToPage, &tx.AccordingToPageLink, &tx.Scrutinizer,
			&tx.ScrutinizerID, &status, &tx.Extinct, &start, &end, &env, &tx.Species,
			&tx.Section, &tx.Subgenus, &tx.Genus, &tx.Subtribe, &tx.Tribe,
			&tx.Subfamily, &tx.Family, &tx.Superfamily, &tx.Suborder, &tx.Order,
			&tx.Subclass, &tx.Class, &tx.Subphylum, &tx.Phylum, &tx.Kingdom,
			&tx.ReferenceID, &tx.Link, &tx.Remarks, &tx.Modified, &tx.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan author row: %w", err)
		}

		tx.TemporalRangeStart = coldp.NewGeoTime(start)
		tx.TemporalRangeEnd = coldp.NewGeoTime(end)

		envs := strings.Split(env, ",")
		envs = gnlib.Map(envs, func(s string) string {
			return strings.TrimSpace(s)
		})
		for _, v := range envs {
			coldpEnv := coldp.NewEnvironment(v)
			if coldpEnv != coldp.UnknownEnv {
				tx.Environment = append(tx.Environment, coldpEnv)
			}
		}

		txStatus := coldp.NewTaxonomicStatus(status)
		if txStatus == coldp.ProvisionallyAcceptedTS {
			tx.Provisional = coldp.ToBool(true)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tx:
			if count%50_000 == 0 {
				util.Progress(count, "taxon")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating name rows: %w", err)
	}

	return nil
}
