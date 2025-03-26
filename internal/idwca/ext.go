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

func (a *idwca) ExtensionSlice(index, offset, limit int) ([][]string, error) {
	if len(a.meta.Extensions) <= index {
		return nil, &arch.ErrExtensionRead{
			Err: fmt.Errorf("extension index is out of bounds: %d", index),
		}
	}
	ext := a.meta.Extensions[index]

	cfg, err := getCsvConfigExt(a.rootDir, *ext)
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

	cfg, err := getCsvConfigExt(a.rootDir, *ext)
	if err != nil {
		return 0, err
	}
	csv := gncsv.New(cfg)
	rowsNum, err := csv.Read(ctx, extCh)
	if err != nil {
		return 0, err
	}

	return rowsNum, nil
}

func getCsvConfigExt(
	rootDir string,
	ext dwca.Extension,
) (config.Config, error) {
	path := filepath.Join(rootDir, ext.Files.Location)
	colSep := ext.FieldsTerminatedBy
	skipHeaders := ext.IgnoreHeaderLines == "1"
	headers := getHeadersExt(ext)
	opts := []config.Option{
		config.OptPath(path),
		config.OptColSep(rune(colSep[0])),
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
	res := make([]string, maxIdx)
	for i := range maxIdx {
		if header, ok := fieldsMap[i]; ok {
			res[i] = header
		} else {
			res[i] = fmt.Sprintf("field%d", i)
		}
	}
	return res
}
