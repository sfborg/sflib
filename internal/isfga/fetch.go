package isfga

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

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
