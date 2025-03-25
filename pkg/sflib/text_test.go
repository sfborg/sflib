package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/sfborg/sflib/pkg/text"
	"github.com/stretchr/testify/assert"
)

func TestNewText(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewText()
	_, ok := sfga.(text.Archive)
	assert.True(ok)
}
