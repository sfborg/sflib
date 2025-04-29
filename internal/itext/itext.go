package itext

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/text"
)

type itext struct {
	filePath   string
	code       nomcode.Code
	jobsNum    int
	parserPool map[nomcode.Code]*sync.Pool
}

func New() text.Archive {
	res := itext{}
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
