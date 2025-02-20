package sfgaio

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/schemaio"
)

// Create copies schema from repo and uses it to create
// new SQLite database.
func (s *sfgaio) Create(dir string) error {
	if s.cfg == nil {
		s.cfg = config.New()
	}
	sch := schemaio.New(s.cfg.GitRepo)
	s.buildDir = dir
	schema, err := sch.Fetch()
	if err != nil {
		slog.Error("Cannot fetch schema", "error", err)
		return err
	}

	schFile := filepath.Join(s.buildDir, "schema.sql")
	err = os.WriteFile(schFile, schema, 0644)
	if err != nil {
		slog.Error("Cannot write schema file", "error", err)
		return err
	}

	s.dbPath, err = makeDb(s.buildDir, schFile)
	if err != nil {
		return err
	}

	return nil
}

// makeDb creates SQLite file from SFGArchive's SQL dump.
func makeDb(workDir, schFile string) (string, error) {
	var err error
	dbPath := filepath.Join(workDir, "schema.sqlite")

	read := fmt.Sprintf(".read %s", schFile)

	cmd := exec.Command("sqlite3", dbPath, read)
	err = cmd.Run()
	if err != nil {
		return "", &sfga.ErrSQLiteLoadSQL{Err: err}
	}
	return dbPath, nil
}
