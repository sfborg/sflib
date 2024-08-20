package sfgaio_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/sfgaio"
	"github.com/stretchr/testify/assert"
)

func TestFetchSchema(t *testing.T) {
	assert := assert.New(t)
	repo := sfga.GitRepo{
		URL: "https://github.com/sfborg/sfga",
		Tag: "v1.2.1",
	}
	tmpPath := filepath.Join(os.TempDir(), "sflibtest")
	fmt.Println(tmpPath)

	s, err := sfgaio.New(repo, tmpPath)
	assert.Nil(err)

	schema, err := s.FetchSchema()
	assert.Nil(err)
	assert.True(len(schema) > 200)
	assert.Contains(string(schema), "CREATE TABLE")
}
