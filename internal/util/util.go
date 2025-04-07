package util

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/pkg/arch"
)

func AssureLocal(src string) (string, string, error) {
	var err error
	var exists bool
	if strings.HasPrefix(src, "http") {
		slog.Info("Downloading from URL", "url", src)
		dlDir, err := os.MkdirTemp("", "sflib-dl")
		if err != nil {
			return dlDir, "", &arch.ErrDownload{URL: src, Err: err}
		}
		src, err = gnsys.Download(src, dlDir, true)
		var dlErr *gnsys.ErrDownload
		if errors.As(err, &dlErr) {
			return dlDir, "", &arch.ErrDownload{URL: dlErr.URL, Err: dlErr.Err}
		}
		if err != nil {
			return dlDir, "", err
		}
		return dlDir, src, nil
	}

	// looks like the file is local, lets check if it exists
	exists, err = gnsys.FileExists(src)
	if err != nil {
		return "", "", err
	}
	if !exists {
		return "", "", &arch.ErrFileNotFound{File: src}
	}

	return "", src, nil
}

func ExtractOrCopy(src, dstDir string) error {
	var err error
	var e gnsys.Extractor
	switch gnsys.GetFileType(src) {
	case gnsys.ZipFT:
		e = gnsys.ExtractZip
	case gnsys.TarFT:
		e = gnsys.ExtractTar
	case gnsys.TarGzFT:
		e = gnsys.ExtractTarGz
	case gnsys.TarBzFT:
		e = gnsys.ExtractTarBz2
	case gnsys.TarXzFt:
		e = gnsys.ExtractTarXz
	default:
		err = copy(src, dstDir)
		if err != nil {
			return err
		}
		return nil
	}
	err = e(src, dstDir)
	if err != nil {
		return &arch.ErrImportArchive{File: src, Err: err}
	}
	return nil
}

func copy(src, dstDir string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstPath := filepath.Join(dstDir, filepath.Base(src))

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func Progress(count int, recordType string) {
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 50))
	fmt.Fprintf(
		os.Stderr,
		"Processed %s %s records\r",
		humanize.Comma(int64(count)),
		recordType,
	)
}

func ProgressEnd() {
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 50))
}
