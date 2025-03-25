package ixsv

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/gnames/gnfmt/gncsv"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/xsv"
)

type ixsv struct {
	filePath   string
	reader     gncsv.Reader
	headers    map[string]int
	code       coldp.NomCode
	jobsNum    int
	parserPool map[coldp.NomCode]*sync.Pool
}

func New() xsv.Archive {
	res := ixsv{}
	return &res
}

func (a *ixsv) Fetch(src, dstDir string) error {
	var err error
	var dlDir string
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
