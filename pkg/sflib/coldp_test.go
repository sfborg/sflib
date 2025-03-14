package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/stretchr/testify/assert"
)

func TestNewCoLDP(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewCoLDP()
	_, ok := sfga.(coldp.Archive)
	assert.True(ok)
}
