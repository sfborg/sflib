package sfgaio_test

import (
	"os"
	"testing"

	"github.com/sfborg/sflib/ent/sfga"
	"github.com/sfborg/sflib/io/sfgaio"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	repo := sfga.GitRepo{
		URL:          "https://github.com/sfborg/sfga",
		Tag:          "v0.3.24",
		ShaSchemaSQL: "b1db9df2e759f",
	}
	tempDir, err := os.MkdirTemp("", "test-create")
	assert.Nil(err)
	defer os.RemoveAll(tempDir)
	s := sfgaio.New()
	err = s.Create(tempDir, repo)
	assert.Nil(err)
	schema = t
}
