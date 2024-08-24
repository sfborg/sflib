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
	var d sfga.DB
	d = dbio.New(dbCache)
	_, ok := d.(sfga.DB)
	assert.True(ok)
}

func TestConnect(t *testing.T) {
	assert := assert.New(t)
	var a sfga.Archive
	var d sfga.DB
	var db *sql.DB
	var err error

	tests := []struct {
		msg   string
		file  string
		isBad bool
	}{
		{"dwca", "dwca.tar.gz", true},
		{"sql", "dinof.sql", false},
		{"sql zip", "dinof.sql.zip", false},
		{"sql tar", "dinof.sql.tar.gz", false},
		{"sqlite", "dinof.sqlite", false},
		{"sqlite zip", "dinof.sqlite.zip", false},
		{"sqlite tar", "dinof.sqlite.tar.gz", false},
	}
	for _, v := range tests {
		sf := filepath.Join("..", "..", "testdata", v.file)
		d = dbio.New(dbCache)

		a, err = archio.New(sf, cache)
		assert.Nil(err)

		err = a.Extract()
		assert.Nil(err)

		db, err = d.Connect()

		if v.isBad {
			assert.NotNil(err)
			assert.Nil(db)
		} else {
			assert.Nil(err)
			assert.NotNil(db)
		}
	}
	err = a.Clean()
	assert.Nil(err)
}

func TestConnectNoExtract(t *testing.T) {
	assert := assert.New(t)

	sf := filepath.Join("..", "..", "testdata", "dwca.tar.gz")
	a, err := archio.New(sf, cache)
	assert.Nil(err)
	err = a.Clean()
	assert.Nil(err)

	d := dbio.New(dbCache)
	db, err := d.Connect()
	assert.NotNil(err)
	assert.Nil(db)
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

	d = dbio.New(dbCache)
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

	d = dbio.New(dbCache)
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
