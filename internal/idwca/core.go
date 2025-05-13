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
	cfg, err := a.getCsvConfigCore()
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
	cfg, err := a.getCsvConfigCore()
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

func (a *idwca) getCsvConfigCore() (config.Config, error) {
	var cfg config.Config
	core := a.meta.Core
	path := filepath.Join(a.rootDir, core.Files.Locations[0])
	var colSep rune
	switch core.FieldsTerminatedBy {
	case "\\t":
		colSep = '\t'
	case "|":
		colSep = '|'
	case ",":
		colSep = ','
	default:
		colSep = 0
		return cfg, fmt.Errorf("Unexpected field separator: '%s'", core.FieldsTerminatedBy)
	}

	var quotes bool
	switch core.FieldsEnclosedBy {
	case "\"":
		quotes = true
	case "":
		quotes = false
	default:
		return cfg, fmt.Errorf("Unexpected fields enclosure: '%s'", core.FieldsEnclosedBy)
	}

	skipHeaders := core.IgnoreHeaderLines == "1"
	headers := getHeadersCore(*core)
	opts := []config.Option{
		config.OptPath(path),
		config.OptColSep(colSep),
		config.OptSkipHeaders(skipHeaders),
		config.OptHeaders(headers),
		config.OptWithQuotes(quotes),
		config.OptBadRowMode(a.cfg.BadRow),
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
