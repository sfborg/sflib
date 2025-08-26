package idwca

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
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

func TestAddTaxonomicStatus(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		msg      string
		nu       coldp.NameUsage
		row      []string
		fieldMap map[string]int
		expected coldp.NameUsage
	}{
		{
			msg: "no status or accepted fields",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "test"},
			fieldMap: map[string]int{
				"id": 0,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				TaxonomicStatus: coldp.UnknownTaxSt,
			},
		},
		{
			msg: "explicit accepted status",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "accepted"},
			fieldMap: map[string]int{
				"id":              0,
				"taxonomicstatus": 1,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				TaxonomicStatus: coldp.AcceptedTS,
			},
		},
		{
			msg: "explicit synonym status",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "synonym"},
			fieldMap: map[string]int{
				"id":              0,
				"taxonomicstatus": 1,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				TaxonomicStatus: coldp.SynonymTS,
			},
		},
		{
			msg: "accepted name usage ID different from ID",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "", "2"},
			fieldMap: map[string]int{
				"id":                  0,
				"taxonomicstatus":     1,
				"acceptednameusageid": 2,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				ParentID:        "2",
				TaxonomicStatus: coldp.SynonymTS,
			},
		},
		{
			msg: "accepted name usage ID same as ID",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "", "1"},
			fieldMap: map[string]int{
				"id":                  0,
				"taxonomicstatus":     1,
				"acceptednameusageid": 2,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				TaxonomicStatus: coldp.AcceptedTS,
			},
		},
		{
			msg: "synonym with explicit status and accepted ID",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "synonym", "2"},
			fieldMap: map[string]int{
				"id":                  0,
				"taxonomicstatus":     1,
				"acceptednameusageid": 2,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				ParentID:        "2",
				TaxonomicStatus: coldp.SynonymTS,
			},
		},
		{
			msg: "only accepted field present with different ID",
			nu:  coldp.NameUsage{ID: "1"},
			row: []string{"1", "2"},
			fieldMap: map[string]int{
				"id":                  0,
				"acceptednameusageid": 1,
			},
			expected: coldp.NameUsage{
				ID:              "1",
				ParentID:        "2",
				TaxonomicStatus: coldp.SynonymTS,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			result := addTaxonomicStatus(tt.nu, tt.row, tt.fieldMap)
			assert.Equal(tt.expected.ID, result.ID)
			assert.Equal(tt.expected.ParentID, result.ParentID)
			assert.Equal(tt.expected.TaxonomicStatus, result.TaxonomicStatus)
		})
	}
}
