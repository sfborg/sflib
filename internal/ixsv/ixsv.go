package ixsv

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/xsv"
)

type ixsv struct {
	cfg        config.Config
	filePath   string
	reader     gncsv.Reader
	headers    map[string]int
	jobsNum    int
	parserPool map[nomcode.Code]chan gnparser.GNparser
}

func New(opts ...config.Option) xsv.Archive {
	cfg := config.New(opts...)
	res := ixsv{cfg: cfg}
	return &res
}

func (a *ixsv) Fetch(src, dstDir string) error {
	var err error
	var dlDir string

	if err = util.AssureEmptyDir(dstDir); err != nil {
		return err
	}

	dlDir, src, err = util.AssureLocal(src)
	if err != nil {
		return err
	}
	if dlDir != "" {
		defer os.RemoveAll(dlDir)
	}

	util.ExtractOrCopy(src, dstDir)
	a.filePath = filepath.Join(dstDir, filepath.Base(src))

	return nil
}

func (a *ixsv) Create(dir string) error {
	return nil
}

func (a *ixsv) Export(outputPath string, isZip bool) error {
	if a.filePath == "" {
		return &arch.ErrFileOpen{Path: "", Err: nil}
	}

	// Add .csv extension if not present
	if filepath.Ext(outputPath) != ".csv" {
		outputPath += ".csv"
	}

	// Copy CSV file to output path if different
	if a.filePath != outputPath {
		if err := copyFile(a.filePath, outputPath); err != nil {
			return err
		}
	}

	slog.Info("CSV file exported", "file", outputPath)

	if isZip {
		if err := createZip(outputPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return &arch.ErrFileOpen{Path: src, Err: err}
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return &arch.ErrFileCreate{File: dst, Err: err}
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return &arch.ErrFileCopy{Src: src, Dst: dst, Err: err}
	}

	return nil
}

func createZip(filePath string) error {
	zipFile := filePath + ".zip"
	f, err := os.Create(zipFile)
	if err != nil {
		return &arch.ErrFileCreate{File: zipFile, Err: err}
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	srcFile, err := os.Open(filePath)
	if err != nil {
		return &arch.ErrFileOpen{Path: filePath, Err: err}
	}
	defer srcFile.Close()

	zipEntry, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}

	_, err = io.Copy(zipEntry, srcFile)
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}

	slog.Info("CSV ZIP file created", "file", zipFile)
	return nil
}

func (a *ixsv) FilePath() string {
	return a.filePath
}

func (a *ixsv) Headers() []string {
	return nil
}

func (a *ixsv) ColdpHeaders() map[string]int {
	return nil
}
