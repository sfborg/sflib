package sflib

import (
	"github.com/sfborg/sflib/internal/ixsv"
	"github.com/sfborg/sflib/pkg/xsv"
)

type xsvArc struct {
	xsv.Archive
}

// NewXsv creates a new xsv.Archive instance, that allow to manage
// an xsv files that contain one scientific name per line.
func NewXsv() xsv.Archive {
	res := xsvArc{
		Archive: ixsv.New(),
	}
	return &res
}
