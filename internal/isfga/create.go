package isfga

import (
	"log/slog"
	"os"
	"path/filepath"
)

// Create copies schema from repo and uses it to create
// new SQLite database.
func (a *isfga) Create(dir string) error {
	sch := NewSchema(a.cfg.GitRepo)
	a.buildDir = dir
	schema, err := sch.Fetch()
	if err != nil {
		slog.Error("Cannot fetch schema", "error", err)
		return err
	}

	schFile := filepath.Join(a.buildDir, "schema.sql")
	err = os.WriteFile(schFile, schema, 0644)
	if err != nil {
		slog.Error("Cannot write schema file", "error", err)
		return err
	}

	a.dbPath, err = makeDb(a.buildDir, schFile)
	if err != nil {
		return err
	}

	return nil
}
