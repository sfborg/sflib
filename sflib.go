// Package sflib is a library that provides functionality to create and manage
// various types of archives for biodiversity data, including text,
// CSV/TSV/PSV, CoLDP, DwCA and SFGA formats.
//
// sflib provides a unified interface for working with different archive
// formats, abstracting away the underlying implementation details. It
// allows users to create new archives of specific types and manage them
// through a consistent set of methods.
//
// The package supports the following archive types:
//
//   - Text: Archives containing plain text files, typically with one
//     scientific name per line.
//   - Xsv: Archives containing delimited files (CSV, TSV, PSV).
//   - Coldp: Archives following the Catalogue of Life Data Package (CoLDP)
//     standard.
//   - DwCA: Archives following the Darwin Core Archive standard.
//   - SFGA: Archives following the Species File Group Archive
//     standard -- SQLite based archive that is close to CoLDP format.
//
// Usage:
//
// To create a new archive, use one of the `New...` functions:
//
//	// Create a new text archive.
//	textArchive := sflib.NewText()
//
//	// Create a new Xsv archive.
//	xsvArchive := sflib.NewXsv()
//
//	// Create a new Dwca archive.
//	dwcaArchive := sflib.NewDwca()
//
//	// Create a new Coldp archive.
//	coldpArchive := sflib.NewColdp()
//
//	// Create a new Sfga archive.
//	sfgaArchive := sflib.NewSfga()
//
// Each of these functions returns an interface type (e.g., `text.Archive`,
// `xsv.Archive`, `dwca.Archive`, `coldp.Archive`, `sfga.Archive`) that can
// be used to interact with the archive.
package sflib

import (
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/icoldp"
	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/internal/itext"
	"github.com/sfborg/sflib/internal/ixsv"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca"
	"github.com/sfborg/sflib/pkg/sfga"
	"github.com/sfborg/sflib/pkg/text"
	"github.com/sfborg/sflib/pkg/xsv"
)

type textArc struct {
	text.Archive
}

type xsvArc struct {
	xsv.Archive
}

type dwcaArc struct {
	dwca.Archive
}

type coldpArc struct {
	coldp.Archive
}

type sfgaArc struct {
	sfga.Archive
}

// NewText creates a new text.Archive instance, which allows managing archives
// containing textual scientific names data. It assumes that a text file has
// UTF-8 encoding and contains one seientific name per line and no other
// information.
//
// Example:
//
//	textArchive := sflib.NewText()
//	// ... use textArchive to manage the archive ...
func NewText() text.Archive {
	res := textArc{
		Archive: itext.New(),
	}
	return &res
}

// NewXsv creates a new xsv.Archive instance, which allows managing archives
// containing xsv (e.g., CSV, TSV or PSV) files.
//
// Example:
//
//	xsvArchive := sflib.NewXsv()
//	// ... use xsvArchive to manage the archive ...
func NewXsv() xsv.Archive {
	res := xsvArc{
		Archive: ixsv.New(),
	}
	return &res
}

func NewDwca(opts ...config.Option) dwca.Archive {
	res := dwcaArc{
		Archive: idwca.New(opts...),
	}
	return &res
}

// NewColdp creates a new coldp.Archive instance, which allows managing archives
// following the Catalogue of Life Data Package (CoLDP) standard.
//
// Example:
//
//	coldpArchive := sflib.NewColdp()
//	// ... use coldpArchive to manage the archive ...
func NewColdp(opts ...config.Option) coldp.Archive {
	res := coldpArc{
		Archive: icoldp.New(opts...),
	}
	return &res
}

// NewSfga creates a new sfga.Archive instance, which allows managing archives
// following the Species File Group Archive. It is based on SQLite database
// and is close to CoLDP format.
//
// Example:
//
//	sfgaArchive := sflib.NewSfga()
//	// ... use sfgaArchive to manage the archive ...
func NewSfga(opts ...config.Option) sfga.Archive {
	res := sfgaArc{
		Archive: isfga.New(opts...),
	}
	return &res
}
