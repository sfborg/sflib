package icoldp_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/icoldp"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

var testDir string

func TestMain(m *testing.M) {
	setupGlobal()
	code := m.Run() // Run all tests
	teardownGlobal()
	os.Exit(code)
}

func setupGlobal() {
	var err error
	testDir, err = os.MkdirTemp("", "coldp-test")
	if err != nil {
		panic(err)
	}
}

func teardownGlobal() {
	var err error
	err = os.RemoveAll(testDir)
	if err != nil {
		panic(err)
	}
}

func TestNameExtract(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
	}{
		{"name-yaml", "name/ptero-yaml.zip"},
		{"usg-yaml", "nameusage/ptero-yaml.zip"},
	}

	for _, v := range tests {
		err := gnsys.CleanDir(testDir)
		assert.Nil(err)
		path := filepath.Join("..", "..", "testdata", "coldp", v.path)
		coldp := icoldp.New()
		err = coldp.Import(path, testDir)
		assert.Nil(err)
		assert.NotNil(coldp)

		assert.Nil(err)
	}
}

func TestNoFile(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	coldp := icoldp.New()
	err = coldp.Import("notfile", testDir)
	assert.IsType(&arch.ErrImportArchive{}, err)
}

func TestNoURL(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	coldp := icoldp.New()
	err = coldp.Import("http://not-a-url", testDir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrDownload{}, err)
}

func TestNotZip(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	path := filepath.Join("..", "..", "testdata", "notzip.zip")
	coldp := icoldp.New()
	err = coldp.Import(path, testDir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrImportArchive{}, err)
}

func TestMeta(t *testing.T) {
	assert := assert.New(t)
	var err error
	var meta *coldp.Meta

	tests := []struct {
		msg, path, creatorGiven, contactOrganization, license string
	}{
		{"name-yaml", "name/ptero-yaml.zip", "Donald", "Species 2000", "CC0"},
		{"usg-yaml", "nameusage/ptero-yaml.zip", "Donald", "Species 2000", "cc0"},
	}

	for _, v := range tests {
		err = gnsys.CleanDir(testDir)
		assert.Nil(err)

		path := filepath.Join("..", "..", "testdata", "coldp", v.path)
		coldp := icoldp.New()
		err = coldp.Import(path, testDir)
		assert.Nil(err)
		err = coldp.DirInfo()
		assert.Nil(err)
		meta, err = coldp.Meta()
		assert.Nil(err)
		assert.Equal(v.creatorGiven, meta.Creators[0].Given)
		assert.Equal(v.contactOrganization, meta.Contact.Organization)
		assert.Equal(v.license, meta.License)
	}
}

func TestName(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gnsys.CleanDir(testDir)
	assert.Nil(err)

	a := icoldp.New()
	path := filepath.Join("..", "..", "testdata", "coldp", "name/ptero-yaml.zip")
	err = a.Import(path, testDir)
	assert.Nil(err)

	err = a.DirInfo()
	assert.Nil(err)

	namePath, ok := a.DataPaths()[coldp.NameDT]
	assert.True(ok)

	ch := make(chan coldp.Name)

	go func() {
		for n := range ch {
			fmt.Printf("%#v\n\n", n)
		}
	}()

	err = coldp.Read(a.Config(), namePath, ch)
	assert.Nil(err)
}
