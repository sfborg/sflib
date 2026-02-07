package ixsv

import (
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
		if err := util.CopyFile(a.filePath, outputPath); err != nil {
			return err
		}
	}

	slog.Info("CSV file exported", "file", outputPath)

	if isZip {
		if err := util.CreateZip(outputPath); err != nil {
			return err
		}
	}

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
