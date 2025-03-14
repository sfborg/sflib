package icoldp

import (
	"errors"
	"os"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/pkg/arch"
)

func (a *icoldp) Import(src, dstDir string) error {
	var err error
	var dlDir string

	a.rootDir = dstDir

	if strings.HasPrefix(src, "http") {
		dlDir, err = os.MkdirTemp("", "coldp-download")
		if err != nil {
			return &arch.ErrDownload{URL: src, Err: err}
		}
		defer os.RemoveAll(dlDir)

		src, err = gnsys.Download(src, dlDir, true)
		if err != nil {
			return &arch.ErrDownload{URL: src, Err: err}
		}
	}

	err = a.extract(src)
	if err != nil {
		return &arch.ErrImportArchive{File: src, Err: err}
	}

	return nil
}

func (a *icoldp) extract(src string) error {
	var err error
	ft := gnsys.GetFileType(src)
	switch ft {
	case gnsys.ZipFT:
		err = gnsys.ExtractZip(src, a.rootDir)
	case gnsys.TarGzFT:
		err = gnsys.ExtractTarGz(src, a.rootDir)
	default:
		err = errors.New("unknown file type")
		return &arch.ErrImportArchive{File: src, Err: err}
	}
	if err != nil {
		return &arch.ErrImportArchive{File: src, Err: err}
	}

	return nil
}
