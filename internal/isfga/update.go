package isfga

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/pkg/sfga"
)

var newDb = "newdb"

type tbl struct {
	name   string
	fields []field
}

type field struct {
	cid       int
	name      string
	typeName  string
	notNull   int
	dfltValue sql.NullString
	pk        int
}

// Update transfers data from an outdated SFGA archive to the most current schema.
// It retrieves the table structure from the source database, attaches the target
// database, and transfers the data table by table.
//
// Parameters:
//   - sfgaEmpty: An sfga.Archive object representing the target database with an
//     empty schema.
//   - withParents: an option that would trigger an attempt to create a
//     parent/child hierarchy out of flat one. In case if flat chierarchy is
//     empty, or parent/child hierarchy already exists, the option is ignored.
//
// Returns:
//   - error: An error if any step in the update process fails, otherwise nil.
func (a *isfga) Update(sfgaEmpty sfga.Archive, withParents bool) error {
	slog.Info("Getting tables information")
	tbls, err := a.getTables()
	if err != nil {
		return err
	}

	if withParents {
		shouldBuild, err := a.shouldAddParents()
		if err != nil {
			return err
		}
		if shouldBuild {
			return a.transferDataWithParents(sfgaEmpty, tbls)
		}
	}

	err = a.transferData(sfgaEmpty.DbPath(), tbls)
	if err != nil {
		return err
	}

	return nil
}

// transferDataWithParents transfers data while building parent-child hierarchy.
// It uses NameUsages to populate name, taxon, and synonym tables with hierarchy,
// then copies all other tables normally.
func (a *isfga) transferDataWithParents(
	sfgaEmpty sfga.Archive, tables []tbl,
) error {
	slog.Info("Transferring data and adding parent IDs")

	// Ensure new database is connected before writing
	if _, err := sfgaEmpty.Connect(); err != nil {
		return err
	}

	if err := a.addParents(sfgaEmpty); err != nil {
		return err
	}

	// These tables are populated already by AddParents
	excludeTables := map[string]bool{
		"name":          true,
		"taxon":         true,
		"synonym":       true,
		"name_relation": true,
	}

	// Transfer remaining tables using the standard method
	if err := a.attachDatabase(sfgaEmpty.DbPath()); err != nil {
		return err
	}
	defer func() {
		if detachErr := a.detachDatabase(); detachErr != nil {
			slog.Error("Failed to detach database", "error", detachErr)
		}
	}()

	for _, tableInfo := range tables {
		if excludeTables[tableInfo.name] {
			continue
		}
		if err := a.transferTableData(tableInfo); err != nil {
			return err
		}
	}

	slog.Info("Finished transferring data with hierarchy")
	return nil
}

func (a *isfga) getTables() ([]tbl, error) {
	var res []tbl

	q := `
SELECT tbl_name
  FROM sqlite_master
  WHERE type = 'table' 
    AND tbl_name not in (
	'sqlite_sequence','version', 'match_type', 'nom_code', 'name_part',
	'gender', 'sex', 'estimate_type', 'distribution_status', 'type_status',
	'nom_rel_type', 'nom_status', 'reference_type', 'taxonomic_status',
	'species_interaction_type', 'taxon_concept_rel_type', 'gazetteer',
	'rank', 'geo_time'
	)
`

	rows, err := a.Db().Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table tbl
		err = rows.Scan(&table.name)
		if err != nil {
			return nil, err
		}
		res = append(res, table)
	}

	for i, v := range res {
		res[i], err = a.getFields(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (a *isfga) getFields(table tbl) (tbl, error) {
	q := `PRAGMA table_info ('%s')`
	rows, err := a.Db().Query(fmt.Sprintf(q, table.name))
	if err != nil {
		return table, err
	}
	defer rows.Close()

	var flds []field
	for rows.Next() {
		var fld field
		err := rows.Scan(
			&fld.cid, &fld.name, &fld.typeName,
			&fld.notNull, &fld.dfltValue, &fld.pk,
		)

		if err != nil {
			return table, err
		}
		flds = append(flds, fld)
	}
	table.fields = flds
	return table, nil
}

func (a *isfga) transferData(dbPath string, tables []tbl) error {
	if err := a.attachDatabase(dbPath); err != nil {
		return err
	}
	defer func() {
		if detachErr := a.detachDatabase(); detachErr != nil {
			slog.Error("Failed to detach database", "error", detachErr)
		}
	}()

	for _, tableInfo := range tables {
		if err := a.transferTableData(tableInfo); err != nil {
			return err
		}
	}

	slog.Info("Finished transferring data")
	return nil
}

func (a *isfga) attachDatabase(dbPath string) error {
	_, err := a.Db().Exec("ATTACH DATABASE '" + dbPath + "' as " + newDb)
	if err != nil {
		slog.Error("Failed to attach database", "dbPath", dbPath, "error", err)
	}
	return err
}

func (a *isfga) detachDatabase() error {
	_, err := a.Db().Exec("DETACH DATABASE " + newDb)
	if err != nil {
		slog.Error("Failed to detach database", "error", err)
	}
	return err
}

func (a *isfga) transferTableData(tableInfo tbl) error {
	tableName := tableInfo.name
	slog.Info("Transferring data", "table", tableName)

	columnNames := gnlib.Map(tableInfo.fields, func(f field) string {
		return f.name
	})

	if len(columnNames) == 0 {
		slog.Warn("No fields in table, skipping", "table", tableName)
		return nil
	}

	selectQ := fmt.Sprintf(
		`SELECT %s FROM %s`,
		strings.Join(columnNames, ","), tableName,
	)

	insertQ := fmt.Sprintf(
		`INSERT INTO %s.%s (%s) %s`,
		newDb, tableName, strings.Join(columnNames, ","), selectQ,
	)

	slog.Debug("Executing query", "query", insertQ)
	_, err := a.Db().Exec(insertQ)
	if err != nil {
		slog.Error("Failed to transfer data", "table", tableName, "error", err)
	}
	return err
}

func (a *isfga) shouldAddParents() (bool, error) {
	// if there is at least one parent, do not add parents.
	var hasParents bool
	err := a.db.QueryRow(`
		SELECT EXISTS (
		  SELECT 1 from taxon
		    WHERE col__parent_id IS NOT NULL AND col__parent_id != ''
		)
	`).Scan(&hasParents)
	if err != nil {
		return false, err
	}
	if hasParents {
		slog.Info("Parent IDs already exist, skipping adding parents")
		return false, nil
	}

	q := `
	SELECT EXISTS (
		SELECT 1 FROM taxon
			WHERE
				col__genus != '' OR col__tribe != '' OR col__family != '' OR
				col__order != '' OR col__class != '' OR col__phylum != '' OR
				col__kingdom != ''
)
`
	var hasFlatHierarchy bool
	err = a.db.QueryRow(q).Scan(&hasFlatHierarchy)
	if err != nil {
		return false, err
	}

	if !hasFlatHierarchy {
		slog.Info("Flat hierarchy does not exist, skipping adding parents")
		return false, nil
	}

	return true, nil
}
