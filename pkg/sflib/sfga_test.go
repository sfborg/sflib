package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/sfga"
	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/stretchr/testify/assert"
)

func TestNewSFGA(t *testing.T) {
	assert := assert.New(t)
	res := sflib.NewSFGA()
	_, ok := res.(sfga.Archive)
	assert.True(ok)
}
