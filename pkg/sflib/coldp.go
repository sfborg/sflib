package sflib

import (
	"github.com/sfborg/sflib/internal/icoldp"
	"github.com/sfborg/sflib/pkg/coldp"
)

type coldpArc struct {
	coldp.Archive
}

func NewCoLDP() coldp.Archive {
	res := coldpArc{
		Archive: icoldp.New(),
	}
	return &res
}
