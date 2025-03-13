package isfga

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/pkg/arch"
)

func (a *isfga) Import(src, dstDir string) error {
	var err error
	var dlDir string
	if strings.HasPrefix(src, "http") {
		dlDir, err = os.MkdirTemp("", "sfga-download")
		if err != nil {
			return &arch.ErrDownload{URL: src, Err: err}
		}
		defer os.RemoveAll(dlDir)

		src, err = gnsys.Download(src, dlDir, true)
		if err != nil {
			return &arch.ErrDownload{URL: src, Err: err}
		}
	}

	err = a.extract(src, dstDir)
	if err != nil {
		return &arch.ErrExtractArchive{File: src, Err: err}
	}

	err = a.setDb(dstDir)
	if err != nil {
		return &arch.ErrSQLiteCreateSQL{File: src, Err: err}
	}

	return nil
}

func (a *isfga) extract(src, dstDir string) error {
	var err error
	ft := gnsys.GetFileType(src)
	switch ft {
	case gnsys.ZipFT:
		err = gnsys.ExtractZip(src, dstDir)
	case gnsys.TarGzFT:
		err = gnsys.ExtractTarGz(src, dstDir)
	case gnsys.SqlFT, gnsys.SqliteFT:
		a.copy(src, dstDir)
	default:
		err = errors.New("unknown file type")
		return &arch.ErrExtractArchive{File: src, Err: err}
	}
	return nil
}

func (a *isfga) copy(src, dstDir string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return &arch.ErrFileOpen{File: src, Err: err}
	}
	defer srcFile.Close()

	base := filepath.Base(src)
	dstPath := filepath.Join(dstDir, base)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return &arch.ErrFileCreate{File: dstPath, Err: err}
	}
	defer dstFile.Close()

	buf := make([]byte, 64*1024) // 32KB buffer size (adjust as needed)
	_, err = io.CopyBuffer(dstFile, srcFile, buf)
	if err != nil {
		return &arch.ErrFileCopy{
			Src: srcFile.Name(),
			Dst: dstFile.Name(),
			Err: err,
		}
	}

	return nil
}

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
		return errors.New("too many database files")
	}

	if len(sql) == 1 {
		sqlFile := filepath.Join(dbDir, sql[0])
		a.dbPath, err = makeDb(dbDir, sqlFile)
		if err != nil {
			return err
		}
		return nil
	} else if len(sql) > 1 {
		return errors.New("too many database dump files")
	}

	return errors.New("no database files")
}

// makeDb creates SQLite file from SFGArchive's SQL dump.
func makeDb(workDir, schFile string) (string, error) {
	var err error
	dbPath := filepath.Join(workDir, "schema.sqlite")

	read := fmt.Sprintf(".read %s", schFile)

	cmd := exec.Command("sqlite3", dbPath, read)
	err = cmd.Run()
	if err != nil {
		return "", &arch.ErrSQLiteLoadSQL{Err: err}
	}
	return dbPath, nil
}
