package isfga

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	atlasschema "ariga.io/atlas/sql/schema"
	atlassqlite "ariga.io/atlas/sql/sqlite"
	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/sfga"
)

// Migrate copies the archive to dstDir, applies Atlas schema migration to
// bring it up to the current schema version, updates the version table, and
// returns a new Archive pointing at the migrated copy.
func (a *isfga) Migrate(dstDir string) (sfga.Archive, error) {
	if v := a.Version(); gnlib.CmpVersion(v, config.RepoMinVersion) < 0 {
		return nil, fmt.Errorf("sfga version %s is below minimum supported %s; "+
			"use an older sflib to migrate first", v, config.RepoMinVersion)
	}

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("creating migration dir: %w", err)
	}

	dstPath := filepath.Join(dstDir, filepath.Base(a.dbPath))
	if err := copyFile(a.dbPath, dstPath); err != nil {
		return nil, fmt.Errorf("copying archive: %w", err)
	}

	ctx := context.Background()

	db, err := sql.Open("sqlite", dstPath)
	if err != nil {
		return nil, fmt.Errorf("opening migrated db: %w", err)
	}
	defer db.Close()

	drv, err := atlassqlite.Open(db)
	if err != nil {
		return nil, fmt.Errorf("opening atlas driver: %w", err)
	}

	var sch sfga.Schema
	if a.cfg.LocalSchemaPath != "" {
		sch = NewSchemaWithLocalPath(a.cfg.GitRepo, a.cfg.LocalSchemaPath)
	} else {
		sch = NewSchema(a.cfg.GitRepo)
	}

	schemaSQL, err := sch.Fetch()
	if err != nil {
		return nil, fmt.Errorf("fetching schema.sql: %w", err)
	}

	desired, err := a.desiredSchemaFromSQL(ctx, schemaSQL)
	if err != nil {
		return nil, fmt.Errorf("building desired schema: %w", err)
	}

	current, err := drv.InspectSchema(ctx, "main", nil)
	if err != nil {
		return nil, fmt.Errorf("inspecting current schema: %w", err)
	}

	changes, err := drv.SchemaDiff(current, desired)
	if err != nil {
		return nil, fmt.Errorf("computing schema diff: %w", err)
	}

	if len(changes) == 0 {
		slog.Info("SFGA schema is already current, no migration needed")
	} else {
		if err := drv.ApplyChanges(ctx, changes); err != nil {
			return nil, fmt.Errorf("applying schema changes: %w", err)
		}
		slog.Info("SFGA schema migrated", "changes", len(changes))
	}

	if _, err := db.ExecContext(ctx, "DELETE FROM version"); err != nil {
		return nil, fmt.Errorf("clearing version table: %w", err)
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO version (sf__id) VALUES (?)", config.SchemaVersion); err != nil {
		return nil, fmt.Errorf("setting version: %w", err)
	}

	result := &isfga{cfg: a.cfg}
	result.SetDb(dstPath)
	if _, err := result.Connect(); err != nil {
		return nil, fmt.Errorf("connecting to migrated archive: %w", err)
	}
	return result, nil
}

// desiredSchemaFromSQL builds an in-memory SQLite database from the provided
// schema SQL and returns the Atlas schema object representing the target state.
func (a *isfga) desiredSchemaFromSQL(
	ctx context.Context,
	schemaSQL []byte,
) (*atlasschema.Schema, error) {
	memDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("opening in-memory db: %w", err)
	}
	defer memDB.Close()

	if _, err := memDB.ExecContext(ctx, string(schemaSQL)); err != nil {
		return nil, fmt.Errorf("applying schema.sql: %w", err)
	}

	memDrv, err := atlassqlite.Open(memDB)
	if err != nil {
		return nil, fmt.Errorf("opening atlas driver on in-memory db: %w", err)
	}

	desired, err := memDrv.InspectSchema(ctx, "main", nil)
	if err != nil {
		return nil, fmt.Errorf("inspecting desired schema: %w", err)
	}
	return desired, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
