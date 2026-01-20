package itext

import (
	"os"
	"path/filepath"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/text"
)

type itext struct {
	cfg        config.Config
	filePath   string
	code       nomcode.Code
	jobsNum    int
	parserPool map[nomcode.Code]chan gnparser.GNparser
}

func New(opts ...config.Option) text.Archive {
	cfg := config.New(opts...)
	res := itext{cfg: cfg}
	return &res
}

func (a *itext) Fetch(src, dstDir string) error {
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

func (a *itext) FilePath() string {
	return a.filePath
}

func (a *itext) Create(dir string) error {
	return nil
}

func (a *itext) Export(outputPath string, isZip bool) error {
	return nil
}
