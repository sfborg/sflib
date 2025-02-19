package schemaio_test

import (
	"testing"

	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/schemaio"
	"github.com/stretchr/testify/assert"
)

func TestFetchSchema(t *testing.T) {
	assert := assert.New(t)
	var err error
	repo := sfga.GitRepo{
		URL:          "https://github.com/sfborg/sfga",
		Tag:          "v0.3.24",
		ShaSchemaSQL: "b1db9df2e759f",
	}
	s := schemaio.New(repo)
	schema, err := s.Fetch()
	assert.Nil(err)

	assert.True(len(schema) > 200)
	assert.Contains(string(schema), "CREATE TABLE")

	// check for matching the hash
	repo.ShaSchemaSQL = "1234567"
	s = schemaio.New(repo)
	schema, err = s.Fetch()
	assert.NotNil(err)
	assert.Nil(schema)
}
