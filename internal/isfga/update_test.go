package isfga_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "update")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/ptero_v0.3.31.sqlite"
	oldSfga := isfga.New()
	err = oldSfga.Fetch(path, dir)
	assert.Nil(err)
	_, err = oldSfga.Connect()
	assert.Nil(err)
	version := oldSfga.Version()

	newSfga := isfga.New()
	err = newSfga.Create(filepath.Join(dir))
	assert.Nil(err)

	oldSfga.Update(newSfga)
	_, err = newSfga.Connect()
	assert.Nil(err)
	assert.NotNil(newSfga.Db())
	var res string
	err = newSfga.Db().QueryRow("select count(*) from taxon").Scan(&res)
	assert.Nil(err)
	assert.Equal("1700", res)

	newVersion := newSfga.Version()
	// 1 means bigger
	assert.Equal(1, gnlib.CmpVersion(newVersion, version))
}
