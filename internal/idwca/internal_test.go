package idwca

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/dwca"
	"github.com/stretchr/testify/assert"
)

// TestMeta tests reading of a Catalogue of Life meta.xml file.
func TestMeta(t *testing.T) {
	assert := assert.New(t)
	var m *dwca.Meta
	var err error
	path := filepath.Join("..", "..", "testdata", "dwca", "meta", "col.xml")

	m, err = getMeta(path)
	assert.Nil(err)
	assert.IsType(m, &dwca.Meta{})
	assert.Equal(m.EMLFile, "eml.xml")
	assert.Equal(m.Core.Encoding, "utf-8")
	assert.Equal(m.Core.Files.Locations[0], "Taxon.tsv")
	assert.Equal(m.Core.FieldsTerminatedBy, "\\t")
	assert.Equal(m.Core.ID.Index, "0")
	assert.Equal(m.Core.Fields[7].Index, "7")
	assert.Equal("http://rs.tdwg.org/dwc/terms/taxonRank", m.Core.Fields[7].Term)
	assert.Equal(3, len(m.Extensions))
	assert.Equal(m.Extensions[0].CoreID.Index, "0")
	assert.Equal("http://rs.gbif.org/terms/1.0/isExtinct", m.Extensions[1].Fields[1].Term)
}

func TestBadMeta(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg, path string
	}{
		{"nofile", "not-a-file"},
		{"badfile", "bad.xml"},
	}

	for _, v := range tests {
		path := filepath.Join("..", "..", "testdata", "dwca", "meta", v.path)
		m, err := getMeta(path)
		assert.NotNil(err, v.msg)

		switch v.msg {
		case "nofile":
			var vErr *arch.ErrFileOpen
			assert.True(errors.As(err, &vErr), v.msg)
		case "decoder":
			var vErr *arch.ErrMetaDecoder
			assert.True(errors.As(err, &vErr), v.msg)
		}

		assert.Nil(m, v.msg)
	}
}

func TestEML(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
	}{
		{"eml", "eml.xml"},
		{"medium", "medium.xml"},
		{"small", "small.xml"},
	}
	var e *dwca.EML
	var err error

	for _, v := range tests {
		path := filepath.Join("..", "..", "testdata", "dwca", "eml", v.path)
		e, err = getEML(path)
		assert.Nil(err)
		assert.NotNil(e, v.msg)
	}
}

func TestBadEML(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
	}{
		{"nofile", "not-a-file"},
		{"badfile", "bad.xml"},
	}

	for _, v := range tests {
		path := filepath.Join("..", "..", "testdata", "dwca", "eml", v.path)
		m, err := getEML(path)
		assert.NotNil(err)

		switch v.msg {
		case "nofile":
			var vErr *arch.ErrFileOpen
			assert.True(errors.As(err, &vErr), v.msg)
		case "decoder":
			var vErr *arch.ErrEmlDecoder
			assert.True(errors.As(err, &vErr), v.msg)
		}

		assert.Nil(m, v.msg)
	}
}
