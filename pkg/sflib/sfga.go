package sflib

import (
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/pkg/arch"
)

type sfga struct {
	arch.SFGA
}

func NewSFGA() arch.SFGA {
	res := sfga{
		SFGA: isfga.New(),
	}
	return &res
}
