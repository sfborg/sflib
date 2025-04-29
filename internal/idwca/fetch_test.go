package idwca_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/dwca/diagn"
	"github.com/stretchr/testify/assert"
)

var testDir string

func TestFetch(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "load")
	err := os.Mkdir(dir, 0777)
	assert.Nil(err)

	tests := []struct {
		msg, path string
		badRow    gnfmt.BadRow
		nameType  diagn.SciNameType
		synType   diagn.SynonymType
		hierType  diagn.HierType
	}{
		{
			"gz",
			"aos-birds.tar.gz",
			gnfmt.ErrorBadRow,
			diagn.SciNameFull,
			diagn.SynNone,
			diagn.HierFlat,
		},
		{
			"vascan",
			"vascan.zip",
			gnfmt.ErrorBadRow,
			diagn.SciNameFull,
			diagn.SynAcceptedID,
			diagn.HierBoth,
		},
		{
			"col",
			"col-mini.zip",
			gnfmt.ErrorBadRow,
			diagn.SciNameFull,
			diagn.SynAcceptedID,
			diagn.HierBoth,
		},
		{
			"pipe",
			"data_pipe.tar.gz",
			gnfmt.ProcessBadRow,
			diagn.SciNameFull,
			diagn.SynHierarchy,
			diagn.HierTree,
		},
	}
	for _, v := range tests {
		err := gnsys.CleanDir(dir)
		assert.Nil(err)

		path := filepath.Join("../../testdata/dwca/", v.path)

		opts := []config.Option{config.OptBadRow(v.badRow)}
		a := idwca.New(opts...)
		err = a.Fetch(path, dir)
		assert.Nil(err, v.msg)
		d := a.Diagnostics()
		assert.Equal(v.nameType, d.SciNameType, v.msg)
		assert.Equal(v.synType, d.SynonymType, v.msg)
		assert.Equal(v.hierType, d.HierType, v.msg)
	}
}

func TestBadLoad(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "badload")
	err := os.Mkdir(dir, 0777)
	assert.Nil(err)

	tests := []struct {
		msg, path string
		err       error
	}{
		{"nometa", "meta_absent.tar.gz", arch.ErrMetaFileNotFound},
		{"manymeta", "meta_dupl.tar.gz", arch.ErrMultipleMetaFiles},
	}
	for _, v := range tests {
		err := gnsys.CleanDir(dir)
		assert.Nil(err)

		path := filepath.Join("../../testdata/dwca/", v.path)

		a := idwca.New()
		err = a.Fetch(path, dir)
		assert.True(errors.Is(err, v.err), v.msg)
	}
}
