package archio_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/archio"
	"github.com/stretchr/testify/assert"
)

var cache = filepath.Join(os.TempDir(), "sflib-test")

func TestNewAndClean(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error
	a, err = archio.New("one", "two")
	assert.Nil(a)
	assert.NotNil(err)

	sf := filepath.Join("..", "..", "testdata", "dinof.sql")
	a, err = archio.New(sf, cache)
	assert.Nil(err)
	assert.NotNil(a)
	err = a.Extract()
	assert.Nil(err)
	status := gnsys.GetDirState(cache)
	assert.Equal(gnsys.DirNotEmpty.String(), status.String())

	err = a.Clean()
	status = gnsys.GetDirState(cache)
	assert.Equal(gnsys.DirEmpty, status)
}

func TestDownload(t *testing.T) {
	if !gnsys.Ping("opendata.globalnames.org:80", 3) {
		return
	}
	assert := assert.New(t)
	var a sfga.Archive
	var err error
	sf := "http://opendata.globalnames.org/sfga/182-gymnodiniales-2018-02-17.sql.zip"
	a, err = archio.New(sf, cache)
	assert.Nil(err)
	err = a.Extract()
	assert.Nil(err)

	dbDir := filepath.Join(cache, "db")
	ents, err := os.ReadDir(dbDir)
	assert.Nil(err)
	assert.Equal(1, len(ents))
	assert.True(strings.HasSuffix(ents[0].Name(), ".sql"))
}

func TestExtract(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var err error
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
		sf := filepath.Join("..", "..", "testdata", v.file)
		a, err = archio.New(sf, cache)
		assert.Nil(err)
		err = a.Extract()
		assert.Nil(err)

		dbDir := filepath.Join(cache, "db")
		ents, err := os.ReadDir(dbDir)
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

		assert.Equal(0, other, v.msg)
		if v.isSql {
			assert.Equal(1, sql, v.msg)
			assert.Equal(0, sqlite, v.msg)
		} else {
			assert.Equal(0, sql, v.msg)
			assert.Equal(1, sqlite, v.msg)
		}
		err = a.Clean()
		assert.Nil(err, v.msg)
	}
}
