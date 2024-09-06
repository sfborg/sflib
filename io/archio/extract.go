package archio

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gnames/gnsys"
)

func (a *archio) extract() error {
	switch filepath.Ext(a.sfgaFilePath) {
	case ".zip":
		return a.extractZip()
	case ".gz":
		return a.extractTarGz()
	case ".sql", ".sqlite":
		return a.copy()
	default:
		base := filepath.Base(a.sfgaFilePath)
		err := fmt.Errorf("unknown extension '%s'", base)
		return err
	}
}

func (a *archio) copy() error {
	src, err := os.Open(a.sfgaFilePath)
	if err != nil {
		return err
	}
	defer src.Close()

	base := filepath.Base(a.sfgaFilePath)
	dstPath := filepath.Join(a.dbPath, base)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	buf := make([]byte, 64*1024) // 32KB buffer size (adjust as needed)
	_, err = io.CopyBuffer(dst, src, buf)
	if err != nil {
		return err
	}

	return nil
}

func (a *archio) extractTarGz() error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(a.sfgaFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gzip reader.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(gzReader)
	return a.untar(tr)
}

func (a *archio) extractZip() error {
	// Open the zip file for reading.
	r, err := zip.OpenReader(a.sfgaFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the full path for the file/directory and ensure its directory exists.
		fpath := filepath.Join(a.dbPath, f.Name)
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// If it's a directory, move on to the next entry.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Open the file within the zip.
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create a file in the filesystem.
		outFile, err := os.OpenFile(
			fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode(),
		)
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the contents of the file from the zip to the new file.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *archio) untar(tarReader *tar.Reader) error {
	var writer *os.File
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Get the individual dbFile from the header.
		dbFile := filepath.Join(a.dbPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory.
			err = os.MkdirAll(dbFile, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// Handle regular file.
			writer, err = os.Create(dbFile)
			if err != nil {
				return err
			}
			io.Copy(writer, tarReader)
			writer.Close()
		default:
			return err
		}
	}
	state := gnsys.GetDirState(a.dbPath)
	if state == gnsys.DirEmpty {
		err := fmt.Errorf("bad tar file '%s'", a.sfgaFilePath)
		return err
	}
	return nil
}
