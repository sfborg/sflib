package dwca

import (
	"context"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca/diagn"
)

type Archive interface {
	arch.Packager
	// Meta returns a pointer to the archive's Meta object.
	Meta() *Meta
	// EML returns a pointer to the archive's provenance EML object.
	EML() *EML
	// Diagnostics returns pointer to diagnostics that did run during
	// loading DwCA data.
	Diagnostics() *diagn.Diagnostics
	//LoadCore reads the content of core file and cnverts it to
	// coldp.NameUsage and coldp.Reference objects.
	LoadCore(
		ctx context.Context,
		ch chan<- coldp.Data,
	) error
	// LoadVernacular reads the content of verncaular names file and converts
	// rows to coldp.Vernacular objects.
	LoadVernacular(
		ctx context.Context,
		idx int,
		ext *Extension,
		ch chan<- []coldp.Vernacular,
	) error
	// LoadDistribution reads the content of distribution file and
	// saves restuls to coldp.Distribution objects.
	LoadDistribution(
		ctx context.Context,
		idx int,
		ext *Extension,
		ch chan<- coldp.Data,
	) error
	// CoreSlice takes offset and number of rows of the core file, and returns
	// the requested data. It returns an error in case if something went wrong.
	CoreSlice(offset, limit int) ([][]string, error)
	// CoreStream fills up a channel with data from the core file. It returns
	// the number of read rows, or an error if something went wrong.
	CoreStream(
		ctx context.Context,
		chCore chan<- []string,
	) (int, error)
	// ExtensionSlice takes index of the extension, offset and number of rows.
	// it returns the corresponding rows from the extension file. In case if
	// something went wrong it returns an error.
	ExtensionSlice(index, offset, limit int) ([][]string, error)
	// ExtensionStream feeds fows of an extension to the provided channel. It
	// returns error if something went wrong.
	ExtensionStream(
		ctx context.Context,
		index int,
		ch chan<- []string,
	) (int, error)
}
