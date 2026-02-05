package dwca

import (
	"context"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca/diagn"
)

// Archive combines reading and writing capabilities for DwCA archives.
type Archive interface {
	arch.Packager
	Reader
	Writer
}

// Reader groups all methods for consuming an existing DwCA.
type Reader interface {
	// Meta returns a pointer to the archive's Meta object.
	Meta() *Meta
	// EML returns a pointer to the archive's provenance EML object.
	EML() *EML
	// Diagnostics returns pointer to diagnostics that did run during
	// loading DwCA data.
	Diagnostics() *diagn.Diagnostics
	// LoadCore reads the content of core file and converts it to
	// coldp.NameUsage and coldp.Reference objects.
	LoadCore(
		ctx context.Context,
		ch chan<- coldp.Data,
	) error
	// LoadVernacular reads the content of vernacular names file and converts
	// rows to coldp.Vernacular objects.
	LoadVernacular(
		ctx context.Context,
		idx int,
		ch chan<- []coldp.Vernacular,
	) error
	// LoadDistribution reads the content of distribution file and
	// saves results to coldp.Distribution objects.
	LoadDistribution(
		ctx context.Context,
		idx int,
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
	// It returns the corresponding rows from the extension file. In case if
	// something went wrong it returns an error.
	ExtensionSlice(index, offset, limit int) ([][]string, error)
	// ExtensionStream feeds rows of an extension to the provided channel. It
	// returns error if something went wrong.
	ExtensionStream(
		ctx context.Context,
		index int,
		ch chan<- []string,
	) (int, error)
}

// Writer groups all methods for creating a new DwCA.
// Every method writes its output to files inside the archive's
// rootDir (set by Create). Export then only packages those
// cached files into a zip.
type Writer interface {
	// WriteMeta marshals the Meta struct to meta.xml on disk.
	WriteMeta(*Meta) error
	// WriteEML converts coldp.Meta to EML and writes eml.xml on disk.
	// If meta is nil or has no title, a synthetic EML is generated.
	WriteEML(*coldp.Meta) error
	// WriteCore writes core taxon rows from a NameUsage channel
	// to Taxon.tsv on disk.
	WriteCore(ctx context.Context, ch <-chan coldp.NameUsage) error
	// WriteVernaculars writes vernacular name extension rows
	// to VernacularName.tsv on disk.
	WriteVernaculars(
		ctx context.Context,
		ch <-chan coldp.Vernacular,
	) error
	// WriteDistributions writes distribution extension rows
	// to Distribution.tsv on disk.
	WriteDistributions(
		ctx context.Context,
		ch <-chan coldp.Distribution,
	) error
}
