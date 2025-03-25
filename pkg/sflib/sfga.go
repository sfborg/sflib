package sflib

import (
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/pkg/sfga"
)

type sfgaArc struct {
	sfga.Archive
}

func NewSfga() sfga.Archive {
	res := sfgaArc{
		Archive: isfga.New(),
	}
	return &res
}
