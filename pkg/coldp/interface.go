package coldp

import (
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/arch"
)

type Archive interface {
	arch.Packager

	// DirInfo finds where data and metadata files are located.
	// It also determines if metadata is provided in JSON or YAML
	// format.
	DirInfo() error

	// DataPaths returns a map where a low-case name of a data type is a key and
	// the path to the corresponding file is the value.
	DataPaths() map[DataType]string

	// Config returns configuration settings of archive.
	Config() config.Config

	// Meta returns coldp.Meta struct. If the struct is empty it populates
	// it with data from meta file first.
	Meta() (*Meta, error)

	// WriteMeta writes Meta data to the given path as JSON.
	WriteMeta(meta *Meta, path string) error
}

type DataLoader interface {
	Load(header, row []string) (DataLoader, error)
}

type DataWriter interface {
	Headers() []string
	Row() []string
}
