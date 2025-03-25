package icoldp

import (
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/coldp"
)

type icoldp struct {
	// cfg is the configuration for CoLDP archive.
	cfg config.ConfigCoLDP

	// rootDir is the path where extracted archive resides.
	rootDir string

	// metaPath is the path to metafile.
	metaPath string

	// metaType can be YAML or JSON
	metaType coldp.FileType

	// meta contains metadata. It is nil if empty.
	meta *coldp.Meta

	// dataPaths contains file paths to dataPaths files.
	// key is low-case type of data, the value is the file path.
	dataPaths map[coldp.DataType]string

	// dataType explains what kind of archive is detected.
	// Can be Name or NameUsage. If neither is found
	// it will generate error during harvest process.
	dataType coldp.ArchiveType
}

func New(opts ...config.OptionCoLDP) coldp.Archive {
	cfg := config.NewColdp(opts...)
	res := icoldp{
		cfg:       cfg,
		dataPaths: make(map[coldp.DataType]string),
	}
	return &res
}

func (a *icoldp) Create(dir string) error {
	return nil
}

func (a *icoldp) Export(output string, isZip bool) error {
	return nil
}

func (a *icoldp) Config() config.ConfigCoLDP {
	return a.cfg
}

func (a *icoldp) DataPaths() map[coldp.DataType]string {
	return a.dataPaths
}
