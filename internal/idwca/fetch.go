package idwca

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/internal/util"
	"github.com/sfborg/sflib/pkg/arch"
)

func (a *idwca) Fetch(srcPath, dstDir string) error {
	var err error
	var rootDir string

	slog.Info("Importing data from DwCA file")
	dlDir, src, err := util.AssureLocal(srcPath)
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

	rootDir, err = getRootDir(dstDir)
	if err != nil {
		return err
	}
	a.rootDir = rootDir

	err = a.load()
	if err != nil {
		return err
	}

	a.diagn, err = a.getDiagnostics()
	if err != nil {
		return err
	}

	return nil
}

// getRootDir determines the directory where the files of DarwinCore archive
// reside.
func getRootDir(path string) (string, error) {
	var dirs []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error,
		) error {
			if err != nil {
				return err // handle the error and possibly abort the Walk
			}

			// Check if the current path is the file we're looking for
			if !info.IsDir() && info.Name() == "meta.xml" {
				dir := filepath.Dir(path) // get the directory of the file
				dirs = append(dirs, dir)  // add it to the slice
			}

			return nil
		})
	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		return "", arch.ErrMetaFileNotFound
	}

	if len(dirs) > 1 {
		return "", arch.ErrMultipleMetaFiles
	}

	return dirs[0], nil
}

func (a *idwca) load() error {
	var err error
	metaPath := filepath.Join(a.rootDir, "meta.xml")
	a.meta, err = getMeta(metaPath)
	if err != nil {
		return err
	}

	// create 'flattened' meta
	a.metaSimple = a.meta.Simplify()

	emlPath := "eml.xml"
	if a.meta.EMLFile != "" {
		emlPath = a.meta.EMLFile
	}
	emlPath = filepath.Join(a.rootDir, emlPath)
	a.eml, err = getEML(emlPath)
	if err != nil {
		var fileOpenErr *arch.ErrFileOpen
		if errors.As(err, &fileOpenErr) {
			slog.Error("Cannot open EML file.", "eml", fileOpenErr.Path, "error", err)
		} else {
			return err
		}
	}
	return nil
}
