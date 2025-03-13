package dwca

import (
	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/pkg/sflib"
)

type dwca struct {
	dwca sflib.DwCA
}

func New() sflib.DwCA {
	res := dwca{
		dwca: idwca.New(),
	}
	return &res
}

func (a *dwca) Create(dir string) error {
	return a.dwca.Create(dir)
}

func (a *dwca) Import(src, dst string) error {
	return a.dwca.Import(src, dst)
}

func (a *dwca) Export(output string, isZip bool) error {
	return a.dwca.Export(output, isZip)
}
