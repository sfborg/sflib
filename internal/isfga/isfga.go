package isfga

import (
	"database/sql"
	"errors"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/sfga"
)

type isfga struct {
	cfg config.Config

	// extractDir is the place where the content of SFGA is extracted to.
	extractDir string

	// buildDir
	buildDir string

	// dbPath is the path to the archive's SQLite file.
	dbPath string

	// connection to SQLite database.
	db *sql.DB
}

func New(opts ...config.Option) sfga.Archive {
	res := isfga{
		cfg: config.New(opts...),
	}
	return &res
}

func (s *isfga) Config() config.Config {
	return s.cfg
}

func (s *isfga) SetDb(path string) {
	_ = s.Close()
	s.dbPath = path
}

func (s *isfga) Connect() (*sql.DB, error) {
	var err error
	var db *sql.DB

	if s.db != nil {
		return s.db, nil
	}

	if s.dbPath == "" {
		err = errors.New("the SQLite path is empty")
		return nil, &arch.ErrSQLiteConnect{Err: err}
	}

	db, err = sql.Open("sqlite", s.dbPath)
	if err != nil {
		return nil, &arch.ErrSQLiteConnect{Err: err}
	}

	// Enable in-memory temporary tables
	_, err = db.Exec("PRAGMA temp_store = MEMORY")
	if err != nil {
		return nil, &arch.ErrSQLitePragma{Err: err}
	}

	// Enable Write-Ahead Logging. Allow many reads and one write concurrently,
	// usually boosts write performance.
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, &arch.ErrSQLitePragma{Err: err}
	}

	s.db = db

	return db, nil
}

func (s *isfga) Db() *sql.DB {
	return s.db
}

func (d *isfga) Ping() bool {
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

func (s *isfga) Close() error {
	var err error
	if s.db != nil {
		err = s.db.Close()
		s.db = nil
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *isfga) DbPath() string {
	return s.dbPath
}

func (s *isfga) Version() string {
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

// IsCompatible checks if the provided version is compatible with the target version of isfga.
// It compares the given version string with the version of the isfga instance using gnlib.CmpVersion.
// If the given version is greater than or equal to the target version, it returns true; otherwise, it returns false.
//
// Parameters:
//   - version: A string representing the version to be checked.
//
// Returns:
//   - bool: true if the provided version is compatible, false otherwise.
func (s *isfga) IsCompatible(version string) bool {
	res := gnlib.CmpVersion(version, s.Version()) != -1
	return res
}
