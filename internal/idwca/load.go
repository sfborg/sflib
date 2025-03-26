package idwca

import (
	"errors"
	"log/slog"
	"path/filepath"

	"github.com/sfborg/sflib/pkg/arch"
)

func (a *idwca) load() error {
	var err error
	metaPath := filepath.Join(a.rootDir, "meta.xml")
	a.meta, err = getMeta(metaPath)
	if err != nil {
		return err
	}

	// create 'flattened' meta
	a.meta.Simplify()

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
