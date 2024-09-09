package archio

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
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
		err := &sfga.ErrUnknownExt{File: base}
		return err
	}
}

func (a *archio) copy() error {
	src, err := os.Open(a.sfgaFilePath)
	if err != nil {
		return &sfga.ErrFileOpen{File: a.sfgaFilePath, Err: err}
	}
	defer src.Close()

	base := filepath.Base(a.sfgaFilePath)
	dstPath := filepath.Join(a.dbPath, base)
	dst, err := os.Create(dstPath)
	if err != nil {
		return &sfga.ErrFileCreate{File: dstPath, Err: err}
	}
	defer dst.Close()

	buf := make([]byte, 64*1024) // 32KB buffer size (adjust as needed)
	_, err = io.CopyBuffer(dst, src, buf)
	if err != nil {
		return &sfga.ErrFileCopy{Src: src.Name(), Dst: dst.Name(), Err: err}
	}

	return nil
}

func (a *archio) extractTarGz() error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(a.sfgaFilePath)
	if err != nil {
		return &sfga.ErrFileOpen{File: a.sfgaFilePath, Err: err}
	}
	defer file.Close()

	// Create a new gzip reader.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return &sfga.ErrTarGzReader{File: a.sfgaFilePath, Err: err}
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
		return &sfga.ErrZipReader{File: a.sfgaFilePath, Err: err}
	}
	defer r.Close()

	for _, f := range r.File {
		// Construct the full path for the file/directory and ensure its directory exists.
		fpath := filepath.Join(a.dbPath, f.Name)
		dir := filepath.Dir(fpath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return &sfga.ErrDirCreate{Dir: dir, Err: err}
		}

		// If it's a directory, move on to the next entry.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Open the file within the zip.
		rc, err := f.Open()
		if err != nil {
			return &sfga.ErrFileOpen{File: f.Name, Err: err}
		}
		defer rc.Close()

		// Create a file in the filesystem.
		outFile, err := os.OpenFile(
			fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode(),
		)
		if err != nil {
			return &sfga.ErrFileOpen{File: outFile.Name(), Err: err}
		}
		defer outFile.Close()

		// Copy the contents of the file from the zip to the new file.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return &sfga.ErrFileCopy{Src: f.Name, Dst: outFile.Name(), Err: err}
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
			return &sfga.ErrTarGzReader{Err: err}
		}

		// Get the individual dbFile from the header.
		dbFile := filepath.Join(a.dbPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory.
			err = os.MkdirAll(dbFile, os.FileMode(header.Mode))
			if err != nil {
				return &sfga.ErrDirCreate{Dir: dbFile, Err: err}
			}
		case tar.TypeReg:
			// Handle regular file.
			writer, err = os.Create(dbFile)
			if err != nil {
				return &sfga.ErrFileCreate{File: dbFile, Err: err}
			}
			io.Copy(writer, tarReader)
			writer.Close()
		default:
			return &sfga.ErrFileCreate{File: dbFile, Err: err}
		}
	}
	state := gnsys.GetDirState(a.dbPath)
	if state == gnsys.DirEmpty {
		err := &sfga.ErrEmptyTar{File: a.sfgaFilePath}
		return err
	}
	return nil
}
