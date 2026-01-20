package xsv

import (
	"context"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
)

// Archive provides methods to operate on an CSV/TSV/PSV (XSV) file
// that contains biodiversity data and with headers that correspond
// to either DarwinCore or CoLDP terms.
type Archive interface {
	arch.Packager

	// Load converts xSV file to coldp.NameUsage objects and feeds them to
	// a channel.
	Load(
		ctx context.Context,
		ch chan<- coldp.NameUsage,
		jobsNum int,
		nomCode nomcode.Code,
	) error

	// FilePath returns the path to the file with scientific names.
	FilePath() string

	// Headers returns the fields detected in the first line of the
	// xsv file.
	Headers() []string

	// ColdpHeaders return a map of a ColDP term in low caps and the
	// corresponding index of the field in the xsv file.
	ColdpHeaders() map[string]int
}
