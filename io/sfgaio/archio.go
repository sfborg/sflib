package sfgaio

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/schemaio"
)

type sfgaio struct {
	// downloadDir is a temporary directory where remote SFGArchive
	// would be downloaded for further processing.
	downloadDir string

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

func (s *sfgaio) Create(dir string, repo sfga.GitRepo) error {
	sch := schemaio.New(repo)
	schema, err := sch.Fetch()
	if err != nil {
		slog.Error("Cannot fetch schema", "error", err)
		return err
	}

	schFile := filepath.Join(dir, "schema.sql")
	err = os.WriteFile(schFile, schema, 0644)
	if err != nil {
		slog.Error("Cannot write schema file", "error", err)
		return err
	}
	return nil
}

func (s *sfgaio) Export(dst string, isZip bool) error {
	return nil
}

func (s *sfgaio) Connect() (*sql.DB, error) {
	return nil, nil
}

func (s *sfgaio) Db() *sql.DB {
	return s.db
}

func (s *sfgaio) Ping() bool {
	return false
}

func (s *sfgaio) Close() error {
	return s.db.Close()
}

func (s *sfgaio) DbPath() string {
	return s.dbPath
}

func (s *sfgaio) Version() string {
	return ""
}
