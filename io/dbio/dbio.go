package dbio

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sfborg/sflib/ent/sfga"
	_ "modernc.org/sqlite"
)

type dbio struct {
	dir     string
	file    string
	isSql   bool
	version string
	db      *sql.DB
}

func New(dbDir string) (sfga.DB, error) {
	var err error
	res := &dbio{dir: dbDir}
	res.file, res.isSql, err = dbFile(dbDir)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *dbio) readFromSQL() error {
	var err error
	dbFile := d.file + "ite"

	read := fmt.Sprintf(".read %s", d.file)
	fmt.Println()

	cmd := exec.Command("sqlite3", dbFile, read)
	err = cmd.Run()
	if err != nil {
		return err
	}

	d.isSql = false
	d.file = dbFile
	return nil
}

func (d *dbio) Connect() (*sql.DB, error) {
	var err error
	var db *sql.DB

	if d.isSql {
		err = d.readFromSQL()
		if err != nil {
			return nil, err
		}
	}

	db, err = sql.Open("sqlite", d.file)
	if err != nil {
		return nil, err
	}

	// Enable in-memory temporary tables
	_, err = db.Exec("PRAGMA temp_store = MEMORY")
	if err != nil {
		return nil, err
	}

	// Enable Write-Ahead Logging. Allow many reads and one write concurrently,
	// usually boosts write performance.
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, err
	}

	d.db = db

	return db, nil
}

func (d *dbio) Close() error {
	return d.db.Close()
}

func (d *dbio) FileDB() string {
	return d.file
}

func (d *dbio) MkSfga(sfgaPath string) error {
	return nil
}

func (d *dbio) Version() string {
	if d.version == "" {
		d.setVersion()
	}
	return d.version
}

func (d *dbio) setVersion() {
	if d.db == nil {
		d.Connect()
	}
	var version string
	res := d.db.QueryRow("SELECT id FROM version LIMIT 1")
	err := res.Scan(&version)
	if err == nil {
		d.version = version
	}
}

// dbFile finds the database file in the extracted data from SFGA.
func dbFile(dbDir string) (string, bool, error) {
	es, err := os.ReadDir(dbDir)
	if err != nil {
		return "", false, err
	}
	if len(es) == 0 || len(es) > 1 {
		err = errors.New("archive should held only one file")
		return "", false, err
	}
	f := es[0].Name()
	ext := filepath.Ext(f)
	path := filepath.Join(dbDir, f)

	switch ext {
	case ".sqlite":
		return path, false, nil
	case ".sql":
		return path, true, nil
	default:
		err = fmt.Errorf("extension should be .sql or .sqlite: %s", f)
		return "", false, err
	}
}
