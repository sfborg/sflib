package idwca

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gnames/gnfmt/gncsv"
	"github.com/gnames/gnfmt/gncsv/config"
	dwca "github.com/sfborg/sflib/pkg/dwca"
)

func (a *idwca) CoreSlice(offset, limit int) ([][]string, error) {
	cfg, err := getCsvConfigCore(a.rootDir, *a.meta.Core)
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

func (a *idwca) CoreStream(
	ctx context.Context,
	coreChan chan<- []string,
) (int, error) {
	cfg, err := getCsvConfigCore(a.rootDir, *a.meta.Core)
	if err != nil {
		return 0, err
	}
	csv := gncsv.New(cfg)
	rowsNum, err := csv.Read(ctx, coreChan)
	if err != nil {
		return 0, err
	}

	return rowsNum, nil
}

func getCsvConfigCore(rootDir string, core dwca.Core) (config.Config, error) {
	path := filepath.Join(rootDir, core.Files.Location)
	colSep := core.FieldsTerminatedBy
	skipHeaders := core.IgnoreHeaderLines == "1"
	headers := getHeadersCore(core)
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

func getHeadersCore(core dwca.Core) []string {
	fieldsMap := make(map[int]string)
	fieldsMap[core.ID.Idx] = "taxonid"
	for _, v := range core.Fields {
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
