package text

import (
	"context"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
)

// Archive is the simplest biodiversity archive.
// It consists of a list of names separated by new lines.
type Archive interface {
	arch.Packager

	// Load reads the prepared archive file line by line and feeds the lines
	// into the channel. It is assumed that each line contains one scientific
	// name. Context allows to cancel the process in case if it receives a Done
	// signal. The nomCode setting helps to interprete names according to
	// Botanical code, if needed. The method parses each name and feeds the
	// results into the channel as a NameUsage object. It returns error if
	// anything goes wrong.
	Load(
		ctx context.Context,
		ch chan<- coldp.NameUsage,
		jobsNum int,
		nomCode nomcode.Code,
	) error

	// FilePath returns the path to the file with scientific names.
	FilePath() string
}
