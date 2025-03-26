package idwca_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/idwca"
	"github.com/sfborg/sflib/pkg/arch"
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
	testDir, err = os.MkdirTemp("", "dwca-test")
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

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "load")
	err := os.Mkdir(dir, 0777)
	assert.Nil(err)

	tests := []struct {
		msg, path string
	}{
		{"gz", "aos-birds.tar.gz"},
		{"noeml", "noeml.zip"},
	}
	for _, v := range tests {
		err := gnsys.CleanDir(dir)
		assert.Nil(err)

		path := filepath.Join("../../testdata/dwca/", v.path)

		a := idwca.New()
		err = a.Fetch(path, dir)
		assert.Nil(err, v.msg)
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
