package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewSFGA()
	_, ok := sfga.(arch.SFGA)
	assert.True(ok)
}
