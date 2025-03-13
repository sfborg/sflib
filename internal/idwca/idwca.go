package idwca

import "github.com/sfborg/sflib/pkg/arch"

type idwca struct{}

func New() arch.DwCA {
	res := idwca{}
	return &res
}

func (a *idwca) Import(src, dst string) error {
	return nil
}
func (a *idwca) Create(dir string) error {
	return nil
}
func (a *idwca) Export(outputPath string, isZip bool) error {
	return nil
}
