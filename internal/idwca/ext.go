package idwca

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnfmt/gncsv/config"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/dwca"
)

func (a *idwca) extByIdx(index int) (*dwca.Extension, error) {
	if len(a.meta.Extensions) <= index {
		return nil, &arch.ErrExtensionRead{
			Err: fmt.Errorf("extension index is out of bounds: %d", index),
		}
	}
	return a.meta.Extensions[index], nil
}

func (a *idwca) ExtensionSlice(index, offset, limit int) ([][]string, error) {
	ext, err := a.extByIdx(index)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(a.rootDir, ext.Files.Locations[0])

	cfg, err := getCsvConfigExt(path, *ext)
	if err != nil {
		return nil, err
	}
	csv := gncsv.New(cfg)

	slice, err := csv.ReadSlice(offset, limit)
	if err != nil {
		return nil, err
	}
	return slice, nil
}

func (a *idwca) ExtensionStream(
	ctx context.Context,
	index int,
	extCh chan<- []string,
) (int, error) {
	if len(a.meta.Extensions) <= index {
		return 0, &arch.ErrExtensionRead{
			Err: fmt.Errorf("extension index is out of bounds: %d", index),
		}
	}
	ext := a.meta.Extensions[index]

	var rowsNum int
	for _, v := range ext.Files.Locations {
		path := filepath.Join(a.rootDir, v)
		cfg, err := getCsvConfigExt(path, *ext)
		if err != nil {
			return 0, err
		}
		cfg.BadRowMode = a.cfg.BadRow
		csv := gncsv.New(cfg)
		num, err := csv.Read(ctx, extCh)
		if err != nil {
			return 0, err
		}
		rowsNum += num
	}
	return rowsNum, nil
}

func getCsvConfigExt(
	path string,
	ext dwca.Extension,
) (config.Config, error) {
	var colSep rune
	switch ext.FieldsTerminatedBy {
	case ",":
		colSep = ','
	case "\\t":
		colSep = '\t'
	case "|":
		colSep = '|'
	}
	skipHeaders := ext.IgnoreHeaderLines == "1"
	headers := getHeadersExt(ext)
	opts := []config.Option{
		config.OptPath(path),
		config.OptColSep(colSep),
		config.OptSkipHeaders(skipHeaders),
		config.OptHeaders(headers),
	}
	cfg, err := config.New(opts...)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func getHeadersExt(ext dwca.Extension) []string {
	fieldsMap := make(map[int]string)
	fieldsMap[ext.CoreID.Idx] = "taxonid"
	for _, v := range ext.Fields {
		header := filepath.Base(v.Term)
		fieldsMap[v.Idx] = strings.ToLower(header)
	}
	var maxIdx int
	for k := range fieldsMap {
		if k > maxIdx {
			maxIdx = k
		}
	}
	res := make([]string, maxIdx+1)
	for i := range maxIdx + 1 {
		if header, ok := fieldsMap[i]; ok {
			res[i] = header
		} else {
			res[i] = fmt.Sprintf("field%d", i)
		}
	}
	return res
}
