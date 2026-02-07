package itext

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/text"
)

type itext struct {
	cfg          config.Config
	sfgaFilePath string
	textFilePath string
	code         nomcode.Code
	jobsNum      int
	parserPool   map[nomcode.Code]chan gnparser.GNparser
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
	a.sfgaFilePath = filepath.Join(dstDir, filepath.Base(src))

	return nil
}

// TODO: can we get rid of it?
func (a *itext) FilePath() string {
	return a.sfgaFilePath
}

func (a *itext) Create(dir string) error {
	err := util.AssureEmptyDir(dir)
	if err != nil {
		return err
	}
	a.textFilePath = filepath.Join(dir, "names.txt")

	return nil
}

func (a *itext) Export(outputPath string, isZip bool) error {
	if a.textFilePath == "" {
		return &arch.ErrFileOpen{Path: "", Err: nil}
	}

	// Add .txt extension if not present
	if filepath.Ext(outputPath) != ".txt" {
		outputPath += ".txt"
	}

	// Copy text file to output path if different
	if a.textFilePath != outputPath {
		if err := util.CopyFile(a.textFilePath, outputPath); err != nil {
			return err
		}
	}

	slog.Info("Text file exported", "file", outputPath)

	if isZip {
		if err := util.CreateZip(outputPath); err != nil {
			return err
		}
	}

	return nil
}

func (a *itext) Write(
	ctx context.Context,
	ch <-chan coldp.NameUsage,
) error {
	f, err := os.Create(a.textFilePath)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", a.textFilePath, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for nu := range ch {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		name := nameString(nu.ScientificName, nu.Authorship)
		if _, err := w.WriteString(name + "\n"); err != nil {
			return fmt.Errorf("cannot write name: %w", err)
		}
	}
	return nil
}

// nameString composes a full scientific name string from the name and
// authorship. If the authorship already appears at the end of the
// scientific name it is not duplicated.
func nameString(sciName, authorship string) string {
	sciName = strings.TrimSpace(sciName)
	authorship = strings.TrimSpace(authorship)
	if authorship == "" || strings.HasSuffix(sciName, authorship) {
		return sciName
	}
	return sciName + " " + authorship
}
