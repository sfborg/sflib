package sfgaio

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
)

func (s *sfgaio) Import(src, dstDir string) error {
	var err error
	var dlDir string
	if strings.HasPrefix(src, "http") {
		dlDir, err = os.MkdirTemp("", "sfga-download")
		if err != nil {
			return &sfga.ErrDownload{URL: src, Err: err}
		}
		defer os.RemoveAll(dlDir)

		src, err = gnsys.Download(src, dlDir, true)
		if err != nil {
			return &sfga.ErrDownload{URL: src, Err: err}
		}
	}

	err = s.extract(src, dstDir)
	if err != nil {
		return &sfga.ErrExtractArchive{File: src, Err: err}
	}

	err = s.setDb(dstDir)
	if err != nil {
		return &sfga.ErrSQLiteCreateSQL{File: src, Err: err}
	}

	return nil
}

func (s *sfgaio) extract(src, dstDir string) error {
	var err error
	ft := gnsys.GetFileType(src)
	switch ft {
	case gnsys.ZipFT:
		err = gnsys.ExtractZip(src, dstDir)
	case gnsys.TarGzFT:
		err = gnsys.ExtractTarGz(src, dstDir)
	case gnsys.SqlFT, gnsys.SqliteFT:
		s.copy(src, dstDir)
	default:
		err = errors.New("unknown file type")
		return &sfga.ErrExtractArchive{File: src, Err: err}
	}
	return nil
}

func (s *sfgaio) copy(src, dstDir string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return &sfga.ErrFileOpen{File: src, Err: err}
	}
	defer srcFile.Close()

	base := filepath.Base(src)
	dstPath := filepath.Join(dstDir, base)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return &sfga.ErrFileCreate{File: dstPath, Err: err}
	}
	defer dstFile.Close()

	buf := make([]byte, 64*1024) // 32KB buffer size (adjust as needed)
	_, err = io.CopyBuffer(dstFile, srcFile, buf)
	if err != nil {
		return &sfga.ErrFileCopy{
			Src: srcFile.Name(),
			Dst: dstFile.Name(),
			Err: err,
		}
	}

	return nil
}

func (s *sfgaio) setDb(dbDir string) error {
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
		s.dbPath = filepath.Join(dbDir, bins[0])
		return nil
	} else if len(bins) > 1 {
		return errors.New("too many database files")
	}

	if len(sql) == 1 {
		sqlFile := filepath.Join(dbDir, sql[0])
		s.dbPath, err = makeDb(dbDir, sqlFile)
		if err != nil {
			return err
		}
		return nil
	} else if len(sql) > 1 {
		return errors.New("too many database dump files")
	}

	return errors.New("no database files")
}
