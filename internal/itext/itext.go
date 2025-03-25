package itext

import (
	"bufio"
	"context"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/text"
)

type itext struct {
	filePath string
}

func New() text.Archive {
	res := itext{}
	return &res
}

func (t *itext) Fetch(src, dstDir string) error {
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
	t.filePath = filepath.Join(dstDir, filepath.Base(src))

	return nil
}

func (t *itext) Load(
	ctx context.Context,
	ch chan<- string) error {
	f, err := os.Open(t.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- scanner.Text():
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

func (t *itext) FilePath() string {
	return t.filePath
}

func (t *itext) Create(dir string) error {
	return nil
}

func (t *itext) Export(outputPath string, isZip bool) error {
	return nil
}
