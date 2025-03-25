package text

import (
	"context"

	"github.com/sfborg/sflib/pkg/arch"
)

// Archive is the simplest biodiversity archive.
// It consists of a list of names separated by new lines.
type Archive interface {
	arch.Packager

	// Load reads the prepared archive file line by line and feeds the
	// lines into the channel. It is assumed that each line contains one
	// scientific name. Context allows to cancel the process in case if it
	// receives a Done signal. The method returns an error if something went
	// wrong.
	Load(ctx context.Context, ch chan<- string) error

	// FilePath returns the path to the file with scientific names.
	FilePath() string
}
