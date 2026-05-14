package isfga_test

import (
	"io"
	"log/slog"
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

var (
	testDir string
)

func TestMain(m *testing.M) {
	setupGlobal()
	code := m.Run() // Run all tests
	teardownGlobal()
	os.Exit(code)
}

func setupGlobal() {
	var err error
	testDir, err = os.MkdirTemp("", "sfga-test")
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func teardownGlobal() {
	var err error
	err = os.RemoveAll(testDir)
	if err != nil {
		panic(err)
	}
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "create")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	s := isfga.New()
	err = s.Create(dir)
	assert.Nil(err)

	exists, _ := gnsys.FileExists(filepath.Join(dir, "schema.sql"))
	assert.True(exists)
	exists, _ = gnsys.FileExists(filepath.Join(dir, "schema.sqlite"))
	assert.True(exists)
	assert.True(s.Ping())
}

func TestNotZip(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "nozip")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := filepath.Join("../../testdata/", "notzip.zip")
	sfga := isfga.New()
	err = sfga.Fetch(path, dir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrImportArchive{}, err)
}

func TestExport(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "export")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	s := isfga.New()
	err = s.Create(dir)
	assert.Nil(err)

	sfgaFile := filepath.Join(dir, "tmp")

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

	tests := []struct {
		name    string
		url     string
		files   int
		suffix  string
	}{
		{
			name:   "sqlite zip",
			url:    "http://opendata.globalnames.org/sfga/tests/0147-vascan-2026-05-14-v37.16.sqlite.zip",
			files:  1,
			suffix: ".sqlite",
		},
		{
			name:   "sql zip",
			url:    "http://opendata.globalnames.org/sfga/tests/0147-vascan-2026-05-14-v37.16.sql.zip",
			files:  2,
			suffix: ".sql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			dir := filepath.Join(testDir, "dl-"+tt.name)
			assert.Nil(os.Mkdir(dir, 0755))
			defer os.RemoveAll(dir)

			a := isfga.New()
			assert.Nil(a.Fetch(tt.url, dir))

			ents, err := os.ReadDir(dir)
			assert.Nil(err)
			assert.Equal(tt.files, len(ents))
			assert.True(strings.HasSuffix(ents[0].Name(), tt.suffix))
		})
	}
}

func TestImport(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error

	dir := filepath.Join(testDir, "import")
	err = os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

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
		err = a.Fetch(src, dir)
		assert.Nil(err)

		ents, err := os.ReadDir(dir)
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
		gnsys.CleanDir(dir)

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

	dir := filepath.Join(testDir, "create")
	err = os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

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
		gnsys.CleanDir(dir)

		a = isfga.New()

		if v.file == "" {
			err = a.Create(dir)
			assert.Nil(err)
		} else {
			err = a.Fetch(src, dir)
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

	dir := filepath.Join(testDir, "ver")
	err = os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

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
		gnsys.CleanDir(dir)

		a = isfga.New()

		if v.file == "" {
			err = a.Create(dir)
			assert.Nil(err)
		} else {
			err = a.Fetch(src, dir)
			assert.Nil(err)
		}
		vers := a.Version()
		assert.NotEmpty(vers)
	}
}
