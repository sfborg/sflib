package isfga

import "github.com/sfborg/sflib/pkg/sflib"

type isfga struct{}

func New() sflib.Packager {
	res := isfga{}
	return &res
}

func (a *isfga) Import(src, dst string) error {
	return nil
}
func (a *isfga) Create(dir string) error {
	return nil
}
func (a *isfga) Export(outputPath string, isZip bool) error {
	return nil
}
