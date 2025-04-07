package isfga

import (
	"context"
	"fmt"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *isfga) LoadTypeMaterials(
	ctx context.Context,
	ch chan<- coldp.TypeMaterial,
) error {
	q := `
SELECT
	col__id, col__source_id, col__name_id, col__citation, col__status_id,
	col__institution_code, col__catalog_number, col__reference_id,
	col__locality, col__country, col__latitude, col__longitude, col__altitude,
	col__host, col__sex_id, col__date, col__collector,
	col__associated_sequences, col__link, col__remarks, col__modified,
	col__modified_by
FROM type_material
`
	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot query SFGA type materials: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var tm coldp.TypeMaterial

		var sex, status string
		err = rows.Scan(
			&tm.ID, &tm.SourceID, &tm.NameID, &tm.Citation, &status,
			&tm.InstitutionCode, &tm.CatalogNumber, &tm.ReferenceID, &tm.Locality,
			&tm.Country, &tm.Latitude, &tm.Longitude, &tm.Altitude, &tm.Host,
			&sex, &tm.Date, &tm.Collector, &tm.AssociatedSequences, &tm.Link,
			&tm.Remarks, &tm.Modified, &tm.ModifiedBy,
		)
		if err != nil {
			return fmt.Errorf("cannot scan type material row: %w", err)
		}

		tm.Sex = coldp.NewSex(sex)
		tm.Status = coldp.NewTypeStatus(status)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- tm:
			if count%50_000 == 0 {
				util.Progress(count, "type material")
			}
		}
	}
	util.ProgressEnd()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating type material rows: %w", err)
	}

	return nil
}
