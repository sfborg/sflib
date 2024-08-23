package dbio_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/archio"
	"github.com/sfborg/sflib/io/dbio"
	"github.com/stretchr/testify/assert"
)

var cache = filepath.Join(os.TempDir(), "sflib-test")
var dbCache = filepath.Join(cache, "db")

func TestNew(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var err error
	tests := []struct {
		msg    string
		file   string
		dbFile string
		isSql  bool
	}{
		{"sql", "dinof.sql", "dinof.sql", true},
		{"sql zip", "dinof.sql.zip", "dinof.sql", true},
		{"sql tar", "dinof.sql.tar.gz", "dinof.sql", true},
		{"bin", "dinof.sqlite", "dinof.sqlite", false},
		{"bin zip", "dinof.sqlite.zip", "dinof.sqlite", false},
		{"bin tar", "dinof.sqlite.tar.gz", "dinof.sqlite", false},
	}

	for _, v := range tests {
		sf := filepath.Join("..", "..", "testdata", v.file)
		a, err = archio.New(sf, cache)
		assert.Nil(err)
		err = a.Extract()
		assert.Nil(err)

		d, err = dbio.New(dbCache)
		assert.Nil(err)
		assert.Equal(v.dbFile, filepath.Base(d.FileDB()))
	}
	err = a.Clean()
	assert.Nil(err)
}

func TestNewBad(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var err error
	sf := filepath.Join("..", "..", "testdata", "dwca.tar.gz")
	a, err = archio.New(sf, cache)
	assert.Nil(err)

	err = a.Extract()
	assert.Nil(err)

	d, err = dbio.New(dbCache)
	assert.NotNil(err)
	assert.Nil(d)

	err = a.Clean()
	assert.Nil(err)
}

func TestConnect(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var db *sql.DB
	var err error

	tests := []struct {
		msg  string
		file string
	}{
		{"sql", "dinof.sql"},
		{"sqlite", "dinof.sqlite"},
	}
	for _, v := range tests {
		sf := filepath.Join("..", "..", "testdata", v.file)
		a, err = archio.New(sf, cache)
		assert.Nil(err)

		err = a.Extract()
		assert.Nil(err)

		d, err = dbio.New(dbCache)
		assert.Nil(err)
		db, err = d.Connect()
		assert.Nil(err)
		assert.NotNil(db)
	}
	err = a.Clean()
	assert.Nil(err)
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var err error

	sf := filepath.Join("..", "..", "testdata", "dinof.sqlite")
	a, err = archio.New(sf, cache)
	assert.Nil(err)

	err = a.Extract()
	assert.Nil(err)

	d, err = dbio.New(dbCache)
	assert.Nil(err)

	ver := d.Version()
	assert.True(strings.HasPrefix(ver, "v"))
}

func TestClose(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var err error
	var db *sql.DB

	sf := filepath.Join("..", "..", "testdata", "dinof.sqlite")
	a, err = archio.New(sf, cache)
	assert.Nil(err)

	err = a.Extract()
	assert.Nil(err)

	d, err = dbio.New(dbCache)
	assert.Nil(err)
	db, err = d.Connect()
	assert.Nil(err)
	assert.NotNil(db)

	var ver string
	err = db.QueryRow("SELECT id from version limit 1").Scan(&ver)
	assert.Nil(err)

	err = d.Close()
	assert.Nil(err)

	err = db.QueryRow("SELECT id from version limit 1").Scan(&ver)
	assert.NotNil(err)
}
