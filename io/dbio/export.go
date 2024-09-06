package dbio

import (
	"archive/zip"
	"io"
	"log/slog"
	"os"
	"os/exec"
)

// Export SQLite database. Take in account output, if the file needs
// to be zipped or not
func (d *dbio) Export(outFile string, isBin, isZip bool) error {
	var err error
	if isBin {
		err = d.dumpBinary(outFile)
	} else {
		err = d.dumpSQL(outFile)
	}

	if isZip {
		createZip(outFile)
	}

	if err != nil {
		return err
	}
	return nil
}

func (d *dbio) dumpBinary(outFile string) error {
	var err error
	cmd := exec.Command("sqlite3", d.file, ".backup "+outFile)

	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (d *dbio) dumpSQL(outFile string) error {
	cmd := exec.Command("sqlite3", d.file, ".dump")
	dumpWriter, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer dumpWriter.Close() // Ensure file gets closed

	cmd.Stdout = dumpWriter // Set command's output to the file

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	slog.Info("SQLite SQL file is created", "file", outFile)

	return nil
}

func createZip(outFile string) error {
	zipFile := outFile + ".zip"
	f, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	w, err := os.Open(outFile)
	if err != nil {
		return err
	}
	defer w.Close()

	fileInfo, err := w.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, w)
	if err != nil {
		return err
	}
	return nil
}
