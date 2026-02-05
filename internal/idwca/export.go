package idwca

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/pkg/arch"
)

func (a *idwca) Export(path string, _ bool) error {
	if filepath.Ext(path) != ".zip" {
		path += ".zip"
	}

	zipfile, err := os.Create(path)
	if err != nil {
		return &arch.ErrFileCreate{File: path, Err: err}
	}
	defer zipfile.Close()

	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()

	err = filepath.Walk(a.rootDir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer file.Close()

		relativePath, err := filepath.Rel(a.rootDir, fpath)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relativePath
		header.Method = zip.Deflate

		zipEntry, err := zipwriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)
		return err
	})

	return err
}
