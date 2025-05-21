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
	col__environment_id, col__species, sf__species_id, col__section,
	sf__section_id, col__subgenus, sf__subgenus_id, col__genus, sf__genus_id,
	col__subtribe, sf__subtribe_id, col__tribe, sf__tribe_id, col__subfamily,
	sf__family_id, col__family, sf__family_id, col__superfamily,
	sf__superfamily_id, col__suborder, sf__suborder_id, col__order, sf__order_id,
	col__subclass, sf__subclass_id, col__class, sf__class_id, col__subphylum,
	sf__subphylum_id, col__phylum, sf__phylum_id, col__kingdom, sf__kingdom_id,
	col__reference_id, col__link, col__remarks, col__modified, col__modified_by
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
			&tx.SpeciesID, &tx.Section, &tx.SectionID, &tx.Subgenus, &tx.SubgenusID,
			&tx.Genus, &tx.GenusID, &tx.Subtribe, &tx.SubtribeID, &tx.Tribe,
			&tx.TribeID, &tx.Subfamily, &tx.SubfamilyID, &tx.Family, &tx.FamilyID,
			&tx.Superfamily, &tx.SuperfamilyID, &tx.Suborder, &tx.SuborderID,
			&tx.Order, &tx.OrderID, &tx.Subclass, &tx.SubclassID, &tx.Class,
			&tx.ClassID, &tx.Subphylum, &tx.SubphylumID, &tx.Phylum, &tx.PhylumID,
			&tx.Kingdom, &tx.KingdomID, &tx.ReferenceID, &tx.Link, &tx.Remarks,
			&tx.Modified, &tx.ModifiedBy,
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
