package sfgaio_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/sfgaio"
	"github.com/stretchr/testify/assert"
)

func TestFetchSchema(t *testing.T) {
	assert := assert.New(t)
	var err error
	repo := sfga.GitRepo{
		URL:          "https://github.com/sfborg/sfga",
		Tag:          "v1.2.1",
		ShaSchemaSQL: "e84cc873",
	}
	tmpPath := filepath.Join(os.TempDir(), repo.ShaSchemaSQL)

	s := sfgaio.New(repo, tmpPath)
	err = s.Clean()
	st := gnsys.GetDirState(s.Path())
	assert.Equal(gnsys.DirAbsent, st)

	schema, err := s.FetchSchema()
	assert.Nil(err)

	st = gnsys.GetDirState(s.Path())
	assert.Equal(gnsys.DirNotEmpty, st)

	assert.True(len(schema) > 200)
	assert.Contains(string(schema), "CREATE TABLE")

	// second time it should take cached schema.
	schema, err = s.FetchSchema()
	assert.Nil(err)

	assert.True(len(schema) > 200)
	assert.Contains(string(schema), "CREATE TABLE")

	// check for matching the hash
	repo.ShaSchemaSQL = "1234567"
	s = sfgaio.New(repo, tmpPath)
	schema, err = s.FetchSchema()
	assert.NotNil(err)
	assert.Nil(schema)

	s.Clean()
	st = gnsys.GetDirState(s.Path())
	assert.Equal(gnsys.DirAbsent, st)
}
