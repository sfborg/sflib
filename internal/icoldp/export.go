package icoldp

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func (a *icoldp) Export(path string, _ bool) error {
	sourceDir := a.rootDir
	// Create the zip file.
	if filepath.Ext(path) != ".zip" {
		path += ".zip"
	}
	zipfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Create the zip writer.
	zipwriter := zip.NewWriter(zipfile)
	defer zipwriter.Close()

	// Walk through the directory and add files to the zip.
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the directory itself.  We only want files.
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Create a zip entry for the file.  We need to determine the
		// relative path within the zip.
		relativePath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		zipEntry, err := zipwriter.Create(relativePath)
		if err != nil {
			return err
		}

		// Copy the file contents to the zip entry.
		_, err = io.Copy(zipEntry, file)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
