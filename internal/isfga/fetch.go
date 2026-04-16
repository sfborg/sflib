package isfga

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/arch"
	_ "modernc.org/sqlite"
)

func (a *isfga) Fetch(src, dstDir string) error {
	dlDir, src, err := util.AssureLocal(src)
	if err != nil {
		return err
	}
	if dlDir != "" {
		defer os.RemoveAll(dlDir)
	}

	err = util.ExtractOrCopy(src, dstDir)
	if err != nil {
		return err
	}

	err = a.setDb(dstDir)
	if err != nil {
		return &arch.ErrSQLiteCreateSQL{File: src, Err: err}
	}

	if gnlib.CmpVersion(a.Version(), config.SchemaVersion) == -1 {
		slog.Warn("SFGA schema outdated, migrating",
			"file_version", a.Version(),
			"target_version", config.SchemaVersion,
		)
		migDir, err := os.MkdirTemp("", "sflib-migrate-")
		if err != nil {
			return fmt.Errorf("creating migration temp dir: %w", err)
		}
		defer os.RemoveAll(migDir)

		migrated, err := a.Migrate(migDir)
		if err != nil {
			return fmt.Errorf("auto-migration failed: %w", err)
		}
		_ = migrated.Close()
		_ = a.Close()

		// Replace the original archive in dstDir with the migrated copy so
		// callers see a single, up-to-date .sqlite file.
		if err := os.Rename(migrated.DbPath(), a.dbPath); err != nil {
			return fmt.Errorf("replacing archive with migrated copy: %w", err)
		}

		if a.cfg.MigrateOutputDir != "" {
			outPath := filepath.Join(a.cfg.MigrateOutputDir, filepath.Base(a.dbPath))
			if err := copyFile(a.dbPath, outPath); err != nil {
				return fmt.Errorf("saving migrated copy: %w", err)
			}
			slog.Info("Migrated SFGA saved", "path", outPath)
		}
	}

	// Close so WAL/SHM sidecar files are checkpointed and removed; callers
	// obtain a live connection via Connect().
	_ = a.Close()
	return nil
}

var (
	ErrTooManyDbFiles  = errors.New("too many database files")
	ErrTooManySqlFiles = errors.New("too many database dump files")
	ErrNoDbFiles       = errors.New("no database files")
)

func (a *isfga) setDb(dbDir string) error {
	es, err := os.ReadDir(dbDir)
	if err != nil {
		return err
	}
	var bins []string
	var sql []string
	for _, v := range es {
		name := v.Name()
		switch {
		case strings.HasSuffix(name, ".sqlite"):
			bins = append(bins, name)
		case strings.HasSuffix(name, ".sql"):
			sql = append(sql, name)
		}
	}

	if len(bins) == 1 {
		a.dbPath = filepath.Join(dbDir, bins[0])
		return nil
	} else if len(bins) > 1 {
		return ErrTooManyDbFiles
	}

	if len(sql) == 1 {
		sqlFile := filepath.Join(dbDir, sql[0])
		a.dbPath, err = makeDb(dbDir, sqlFile)
		if err != nil {
			return err
		}
		return nil
	} else if len(sql) > 1 {
		return ErrTooManySqlFiles
	}

	return ErrNoDbFiles
}

// makeDb creates SQLite file from SFGArchive's SQL dump.
func makeDb(workDir, schFile string) (string, error) {
	dbPath := filepath.Join(workDir, "schema.sqlite")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return "", &arch.ErrSQLiteConnect{Err: err}
	}
	defer db.Close()

	schema, err := os.ReadFile(schFile)
	if err != nil {
		return "", &arch.ErrFileOpen{Path: schFile, Err: err}
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return "", &arch.ErrSQLiteExec{Err: err}
	}
	return dbPath, nil
}
