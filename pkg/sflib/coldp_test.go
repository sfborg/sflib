package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/stretchr/testify/assert"
)

func TestNewColdp(t *testing.T) {
	assert := assert.New(t)
	res := sflib.NewColdp()
	_, ok := res.(coldp.Archive)
	assert.True(ok)
}
