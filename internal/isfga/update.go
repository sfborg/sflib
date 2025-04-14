package isfga

import (
	"database/sql"
	"fmt"
	"log/slog"

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
//
// Returns:
//   - error: An error if any step in the update process fails, otherwise nil.
func (a *isfga) Update(sfgaEmpty sfga.Archive) error {
	slog.Info("Getting tables information")
	tbls, err := a.getTables()
	if err != nil {
		return err
	}

	err = a.transferData(sfgaEmpty.DbPath(), tbls)
	if err != nil {
		return err
	}

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
		joinWithComma(columnNames), tableName,
	)

	insertQ := fmt.Sprintf(
		`INSERT INTO %s.%s (%s) %s`,
		newDb, tableName, joinWithComma(columnNames), selectQ,
	)

	slog.Debug("Executing query", "query", insertQ)
	_, err := a.Db().Exec(insertQ)
	if err != nil {
		slog.Error("Failed to transfer data", "table", tableName, "error", err)
	}
	return err
}

func joinWithComma(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, str := range strs[1:] {
		result += "," + str
	}
	return result
}
