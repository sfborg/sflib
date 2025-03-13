package icoldp

import "github.com/sfborg/sflib/pkg/arch"

type icoldp struct{}

func New() arch.CoLDP {
	res := icoldp{}
	return &res
}

func (a *icoldp) Create(dir string) error {
	return nil
}

func (a *icoldp) Import(src, dst string) error {
	return nil
}

func (a *icoldp) Export(output string, isZip bool) error {
	return nil
}
