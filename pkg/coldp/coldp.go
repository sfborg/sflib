package coldp

import (
	"github.com/sfborg/sflib/internal/icoldp"
	"github.com/sfborg/sflib/pkg/sflib"
)

type coldp struct {
	coldp sflib.CoLDP
}

func New() sflib.CoLDP {
	res := coldp{
		coldp: icoldp.New(),
	}
	return &res
}

func (a *coldp) Create(dir string) error {
	return a.coldp.Create(dir)
}

func (a *coldp) Import(src, dst string) error {
	return a.coldp.Import(src, dst)
}

func (a *coldp) Export(output string, isZip bool) error {
	return a.coldp.Export(output, isZip)
}
