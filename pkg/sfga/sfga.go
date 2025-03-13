package sfga

import (
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/pkg/sflib"
)

type sfga struct {
	sfga sflib.SFGA
}

func New() sflib.SFGA {
	res := sfga{
		sfga: isfga.New(),
	}
	return &res
}

func (a *sfga) Create(dir string) error {
	return a.sfga.Create(dir)
}

func (a *sfga) Import(src, dst string) error {
	return a.sfga.Import(src, dst)
}

func (a *sfga) Export(output string, isZip bool) error {
	return a.sfga.Export(output, isZip)
}
