package isfga

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"os/exec"

	"github.com/sfborg/sflib/pkg/arch"
)

// Export SQLite database. Take in account output, if the file needs
// to be zipped or not
func (a *isfga) Export(outFile string, isZip bool) error {
	var err error
	err = a.dumpSQL(outFile + ".sql")
	if err != nil {
		return err
	}
	err = a.dumpBinary(outFile + ".sqlite")
	if err != nil {
		return err
	}

	if isZip {
		err = createZip(outFile + ".sql")
		if err != nil {
			return err
		}
		err = createZip(outFile + ".sqlite")
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *isfga) dumpBinary(outFile string) error {
	var err error
	cmd := exec.Command("sqlite3", a.dbPath, ".backup "+outFile)

	if err = cmd.Run(); err != nil {
		return &arch.ErrSQLiteCreateBinary{File: outFile, Err: err}
	}
	slog.Info("SQLite binary file is created", "file", outFile)

	return nil
}

func (a *isfga) dumpSQL(outFile string) error {
	cmd := exec.Command("sqlite3", a.dbPath, ".dump")
	dumpWriter, err := os.Create(outFile)
	if err != nil {
		return &arch.ErrSQLiteCreateSQL{File: outFile, Err: err}
	}
	defer dumpWriter.Close() // Ensure file gets closed

	cmd.Stdout = dumpWriter // Set command's output to the file

	if err = cmd.Start(); err != nil {
		return &arch.ErrFileCopy{Src: a.dbPath, Dst: outFile, Err: err}
	}

	if err = cmd.Wait(); err != nil {
		return &arch.ErrFileCopy{Src: a.dbPath, Dst: outFile, Err: err}
	}

	slog.Info("SQLite SQL file is created", "file", outFile)

	return nil
}

func createZip(outFile string) error {
	zipFile := outFile + ".zip"
	f, err := os.Create(zipFile)
	if err != nil {
		return &arch.ErrFileCreate{File: zipFile, Err: err}
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	w, err := os.Open(outFile)
	if err != nil {
		return &arch.ErrFileOpen{Path: outFile, Err: err}
	}
	defer w.Close()

	fileInfo, err := w.Stat()
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}

	_, err = io.Copy(writer, w)
	if err != nil {
		return &arch.ErrZipCreate{File: zipFile, Err: err}
	}
	slog.Info("SQLite ZIP file is created", "file", zipFile)
	return nil
}
