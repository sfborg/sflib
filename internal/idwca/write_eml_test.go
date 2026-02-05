package idwca_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestWriteEML(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "write-eml-test")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	a := idwca.New()
	err = a.Create(dir)
	assert.Nil(err)

	// nil Meta should produce synthetic file
	err = a.WriteEML(nil)
	assert.Nil(err)

	bs, err := os.ReadFile(filepath.Join(dir, "eml.xml"))
	assert.Nil(err)
	content := string(bs)
	assert.True(strings.Contains(content, "Placeholder DwCA"))
	assert.True(strings.Contains(content, "auto-generated"))
	assert.True(strings.Contains(content, `eml:eml`))
	assert.True(strings.Contains(content, `xmlns:eml=`))
	assert.Equal("Placeholder DwCA (no metadata provided)", a.EML().Dataset.Title)

	// empty-title Meta should also produce synthetic file
	err = a.WriteEML(&coldp.Meta{})
	assert.Nil(err)
	bs, err = os.ReadFile(filepath.Join(dir, "eml.xml"))
	assert.Nil(err)
	assert.True(strings.Contains(string(bs), "Placeholder DwCA"))

	// real Meta should be converted to EML and written
	meta := &coldp.Meta{
		Key:         "test-dataset-1",
		Title:       "Test Dataset",
		Description: "A test abstract.",
		Issued:      "2024-01-01",
		DOI:         "10.1234/test",
		License:     "CC-BY-4.0",
		Keywords:    []string{"taxonomy", "test"},
		Creators: []coldp.Actor{
			{
				Given:  "Jane",
				Family: "Doe",
				Email:  "jane@example.com",
			},
		},
		Contact: &coldp.Actor{
			Given:        "John",
			Family:       "Smith",
			Organization: "Test Org",
			Country:      "USA",
		},
		Editors: []coldp.Actor{
			{Given: "Editor", Family: "One"},
		},
		Contributors: []coldp.Actor{
			{Given: "Contrib", Family: "Two"},
		},
		GeographicScope: "Global",
	}
	err = a.WriteEML(meta)
	assert.Nil(err)

	bs, err = os.ReadFile(filepath.Join(dir, "eml.xml"))
	assert.Nil(err)
	content = string(bs)
	assert.True(strings.Contains(content, "Test Dataset"))
	assert.True(strings.Contains(content, "Jane"))
	assert.True(strings.Contains(content, "Doe"))
	assert.True(strings.Contains(content, "A test abstract."))
	assert.True(strings.Contains(content, `packageId="test-dataset-1"`))
	assert.True(strings.Contains(content, `<dataset`))
	assert.True(strings.Contains(content, `id="test-dataset-1"`))
	assert.True(strings.Contains(content, "10.1234/test"))
	assert.True(strings.Contains(content, "CC-BY-4.0"))
	assert.True(strings.Contains(content, "John"))
	assert.True(strings.Contains(content, "Smith"))
	assert.True(strings.Contains(content, "Editor"))
	assert.True(strings.Contains(content, "Contrib"))
	assert.True(strings.Contains(content, "Global"))

	// Verify EML struct was populated correctly via Reader interface
	eml := a.EML()
	assert.Equal("Test Dataset", eml.Dataset.Title)
	assert.Equal(1, len(eml.Dataset.Creators))
	assert.Equal("Jane", eml.Dataset.Creators[0].IndividualName.GivenName)
	assert.Equal(1, len(eml.Dataset.Contacts))
	assert.Equal(2, len(eml.Dataset.AssociatedParties))
}
