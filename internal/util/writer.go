package util

import (
	"bufio"
	"encoding/csv"
	"os"
)

type writer struct {
	file           *os.File
	path           string
	bufferedWriter *bufio.Writer
	*csv.Writer
	Count int
}

func NewWriter(path string, colSep rune) (*writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	bufferedWriter := bufio.NewWriter(f)
	csvWriter := csv.NewWriter(bufferedWriter)
	csvWriter.Comma = colSep
	csvWriter.UseCRLF = false

	res := writer{
		path:           path,
		file:           f,
		bufferedWriter: bufferedWriter,
		Writer:         csvWriter,
	}
	return &res, nil
}

func (w *writer) Close() error {
	var err error
	w.Flush()
	if err = w.Error(); err != nil {
		return err
	}
	w.bufferedWriter.Flush()
	err = w.file.Close()
	if err != nil {
		return err
	}
	// remove the file if it is empty
	if w.Count == 0 {
		if err := os.Remove(w.path); err != nil {
			return err
		}
	}
	return nil
}
