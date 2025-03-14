package isfga_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/sfga"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	s := isfga.New()
	err = s.Create(tempDir)
	assert.Nil(err)

	exists, _ := gnsys.FileExists(filepath.Join(tempDir, "schema.sql"))
	assert.True(exists)
	exists, _ = gnsys.FileExists(filepath.Join(tempDir, "schema.sqlite"))
	assert.True(exists)
	assert.True(s.Ping())
}

func TestNotZip(t *testing.T) {
	assert := assert.New(t)
	tmpDir, err := os.MkdirTemp("", "coldp-test")
	assert.Nil(err)
	defer os.RemoveAll(tmpDir)

	path := filepath.Join("..", "..", "testdata", "notzip.zip")
	coldp := isfga.New()
	err = coldp.Import(path, tmpDir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrImportArchive{}, err)
}

func TestExport(t *testing.T) {
	assert := assert.New(t)
	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	s := isfga.New()
	err = s.Create(tempDir)
	assert.Nil(err)

	sfgaFile := filepath.Join(tempDir, "tmp")

	err = s.Export(sfgaFile, true)
	assert.Nil(err)

	file := sfgaFile + ".sql"
	exists, _ := gnsys.FileExists(file)
	assert.True(exists)

	file = sfgaFile + ".sqlite"
	exists, _ = gnsys.FileExists(file)
	assert.True(exists)

	file = sfgaFile + ".sql.zip"
	exists, _ = gnsys.FileExists(file)
	assert.True(exists)

	file = sfgaFile + ".sqlite.zip"
	exists, _ = gnsys.FileExists(file)
	assert.True(exists)
}

func TestDownload(t *testing.T) {
	if !gnsys.Ping("opendata.globalnames.org:80", 3) {
		return
	}
	assert := assert.New(t)

	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	var a sfga.Archive
	sf := "http://opendata.globalnames.org/sfga/147-vascan-2025-01-31.sql.zip"
	a = isfga.New()
	assert.Nil(err)
	err = a.Import(sf, tempDir)
	assert.Nil(err)

	ents, err := os.ReadDir(tempDir)
	assert.Nil(err)
	// sqlite should be created from sql, therefore 2 fiels
	assert.Equal(2, len(ents))
	assert.True(strings.HasSuffix(ents[0].Name(), ".sql"))
}

func TestImport(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error

	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		msg    string
		file   string
		isSql  bool
		isArch bool
	}{
		{"sql", "dinof.sql", true, false},
		{"sql zip", "dinof.sql.zip", true, true},
		{"sql tar", "dinof.sql.tar.gz", true, true},
		{"bin", "dinof.sqlite", false, false},
		{"bin zip", "dinof.sqlite.zip", false, true},
		{"bin tar", "dinof.sqlite.tar.gz", false, true},
	}

	for _, v := range tests {
		src := filepath.Join("..", "..", "testdata", "sfga", v.file)
		a = isfga.New()
		err = a.Import(src, tempDir)
		assert.Nil(err)

		ents, err := os.ReadDir(tempDir)
		assert.Nil(err)

		var sql, sqlite, other int
		for _, e := range ents {
			switch filepath.Ext(e.Name()) {
			case ".sql":
				sql++
			case ".sqlite":
				sqlite++
			default:
				other++
			}
		}
		gnsys.CleanDir(tempDir)

		assert.Equal(0, other, v.msg)
		if v.isSql {
			assert.Equal(1, sql, v.msg)
			// sqlite is created during import
			assert.Equal(1, sqlite, v.msg)
		} else {
			assert.Equal(0, sql, v.msg)
			assert.Equal(1, sqlite, v.msg)
		}
		assert.Nil(err, v.msg)
	}
}

func TestConnect(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error

	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		msg  string
		file string
	}{
		{"sql", "dinof.sql"},
		{"sql zip", "dinof.sql.zip"},
		{"sql tar", "dinof.sql.tar.gz"},
		{"bin", "dinof.sqlite"},
		{"bin zip", "dinof.sqlite.zip"},
		{"bin tar", "dinof.sqlite.tar.gz"},
		{"new", ""},
	}

	for _, v := range tests {
		src := filepath.Join("..", "..", "testdata", "sfga", v.file)
		gnsys.CleanDir(tempDir)

		a = isfga.New()

		if v.file == "" {
			err = a.Create(tempDir)
			assert.Nil(err)
		} else {
			err = a.Import(src, tempDir)
			assert.Nil(err)
		}

		_, err = a.Connect()
		assert.True(a.Ping())
	}
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error

	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		msg  string
		file string
	}{
		{"sql", "dinof.sql"},
		{"sql zip", "dinof.sql.zip"},
		{"sql tar", "dinof.sql.tar.gz"},
		{"bin", "dinof.sqlite"},
		{"bin zip", "dinof.sqlite.zip"},
		{"bin tar", "dinof.sqlite.tar.gz"},
		{"new", ""},
	}

	for _, v := range tests {
		src := filepath.Join("..", "..", "testdata", "sfga", v.file)
		gnsys.CleanDir(tempDir)

		a = isfga.New()

		if v.file == "" {
			err = a.Create(tempDir)
			assert.Nil(err)
		} else {
			err = a.Import(src, tempDir)
			assert.Nil(err)
		}
		vers := a.Version()
		assert.NotEmpty(vers)
	}
}
