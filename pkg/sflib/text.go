package sflib

import (
	"github.com/sfborg/sflib/internal/itext"
	"github.com/sfborg/sflib/pkg/text"
)

type textArc struct {
	text.Archive
}

// NewText creates a new text.Archive instance, that allow to manage
// an text files that contain one scientific name per line.
func NewText() text.Archive {
	res := textArc{
		Archive: itext.New(),
	}
	return &res
}
