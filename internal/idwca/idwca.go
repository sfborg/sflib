package idwca

import (
	"sync"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/dwca"
	"github.com/sfborg/sflib/pkg/dwca/diagn"
	"github.com/sfborg/sflib/pkg/parser"
)

type idwca struct {
	cfg config.Config
	// rootDir is the directory where meta.xml file is resided.
	// it is not always the 'root' directory of extracted archive.
	rootDir string
	// meta contains data from meta.xml file: name and location of files,
	// name and position of fields in csv files etc.
	meta *dwca.Meta
	// metaSimple is a simplified version of Meta. It allows working
	// with metadata easier.
	metaSimple *dwca.MetaSimple
	// eml provides data about the archive and its provenance.
	eml *dwca.EML
	// diagn provides types of ScientificName, Hieararchy, Synonymy.
	diagn *diagn.Diagnostics
	// parserPool contains parsers for names in botanical and zoological codes.
	parserPool map[nomcode.Code]*sync.Pool
}

func New(opts ...config.Option) dwca.Archive {
	cfg := config.New(opts...)
	res := idwca{
		cfg:        cfg,
		parserPool: parser.Pool(cfg.JobsNum),
	}
	return &res
}

func (a *idwca) Create(dir string) error {
	// TODO implement
	return nil
}
func (a *idwca) Export(outputPath string, isZip bool) error {
	// TODO implement
	return nil
}

func (a *idwca) EML() *dwca.EML {
	return a.eml
}

func (a *idwca) Meta() *dwca.Meta {
	return a.meta
}

func (a *idwca) Diagnostics() *diagn.Diagnostics {
	return a.diagn
}
