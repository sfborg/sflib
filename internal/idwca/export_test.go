package idwca_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "export-test")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	srcDir := filepath.Join(dir, "src")
	err = os.Mkdir(srcDir, 0755)
	assert.Nil(err)

	a := idwca.New()
	err = a.Create(srcDir)
	assert.Nil(err)

	// Write an EML so there's a file to package
	err = a.WriteEML(&coldp.Meta{
		Title:       "Export Test",
		Description: "Testing export.",
	})
	assert.Nil(err)

	// Export to zip
	outPath := filepath.Join(dir, "test.zip")
	err = a.Export(outPath, true)
	assert.Nil(err)

	// Verify zip contents
	r, err := zip.OpenReader(outPath)
	assert.Nil(err)
	defer r.Close()

	names := make(map[string]bool)
	for _, f := range r.File {
		names[f.Name] = true
	}
	assert.True(names["eml.xml"])

	// Export without .zip extension should append it
	outPath2 := filepath.Join(dir, "test2")
	err = a.Export(outPath2, true)
	assert.Nil(err)

	_, err = os.Stat(outPath2 + ".zip")
	assert.Nil(err)
}
