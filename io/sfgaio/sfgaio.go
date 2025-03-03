package sfgaio

import (
	"database/sql"
	"errors"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/ent/sfga"
)

type sfgaio struct {
	cfg *config.Config
	// extractDir is the place where the content of SFGA is extracted to.
	extractDir string

	// buildDir
	buildDir string

	// dbPath is the path to the archive's SQLite file.
	dbPath string

	// connection to SQLite database.
	db *sql.DB
}

// New creates an empty instance
func New() sfga.Archive {
	return &sfgaio{}
}

func (s *sfgaio) Connect() (*sql.DB, error) {
	var err error
	var db *sql.DB

	if s.dbPath == "" {
		err = errors.New("the SQLite path is empty")
		return nil, &sfga.ErrSQLiteConnect{Err: err}
	}

	db, err = sql.Open("sqlite", s.dbPath)
	if err != nil {
		return nil, &sfga.ErrSQLiteConnect{Err: err}
	}

	// Enable in-memory temporary tables
	_, err = db.Exec("PRAGMA temp_store = MEMORY")
	if err != nil {
		return nil, &sfga.ErrSQLitePragma{Err: err}
	}

	// Enable Write-Ahead Logging. Allow many reads and one write concurrently,
	// usually boosts write performance.
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, &sfga.ErrSQLitePragma{Err: err}
	}

	s.db = db

	return db, nil
}

func (s *sfgaio) Db() *sql.DB {
	return s.db
}

func (d *sfgaio) Ping() bool {
	var err error
	if d.db == nil {
		_, err = d.Connect()
		if err != nil {
			return false
		}
	}
	err = d.db.Ping()
	if err != nil {
		return false
	}
	return true
}

func (s *sfgaio) Close() error {
	return s.db.Close()
}

func (s *sfgaio) DbPath() string {
	return s.dbPath
}

func (s *sfgaio) Version() string {
	var err error
	if s.db == nil {
		_, err = s.Connect()
		if err != nil {
			return ""
		}
	}

	var version string
	err = s.db.QueryRow("SELECT id FROM version LIMIT 1").Scan(&version)
	if err != nil {
		return ""
	}
	return version
}

// IsCompatible checks if the provided version is compatible with the current version of sfgaio.
// It compares the given version string with the version of the sfgaio instance using gnlib.CmpVersion.
// If the given version is greater than or equal to the current version, it returns true; otherwise, it returns false.
//
// Parameters:
//   - version: A string representing the version to be checked.
//
// Returns:
//   - bool: true if the provided version is compatible, false otherwise.
func (s *sfgaio) IsCompatible(version string) bool {
	res := gnlib.CmpVersion(version, s.Version()) != -1
	return res
}
