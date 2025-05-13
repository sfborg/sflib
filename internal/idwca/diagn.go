package idwca

import (
	"errors"
	"strings"

	"github.com/sfborg/sflib/pkg/dwca/diagn"
)

// getDiagnostics provides an idea how DwCA organizes scientific names,
// synonyms and hieararchy.
func (a *idwca) getDiagnostics() (*diagn.Diagnostics, error) {
	cs, exts, err := a.coreSample()
	if err != nil {
		return nil, err
	}
	if cs == nil {
		return nil, errors.New("no data in the core file")
	}

	res := diagn.New(cs, exts)
	return res, nil
}

func (a *idwca) coreSample() (
	[]map[string]string,
	map[string]string,
	error,
) {
	dt, err := a.CoreSlice(0, 1000)
	if err != nil {
		return nil, nil, err
	}
	m := a.metaSimple

	exts := make(map[string]string)
	for k, v := range m.ExtensionsData {
		exts[k] = strings.ToLower(v.Locations[0])
	}
	coreRows := make([]map[string]string, len(dt))
	for i, row := range dt {
		coreRows[i] = make(map[string]string)
		for j, val := range row {
			if m.CoreData.FieldsIdx[j].Term == "" {
				continue
			}
			coreRows[i][m.CoreData.FieldsIdx[j].Term] = val
		}
	}
	return coreRows, exts, nil
}
