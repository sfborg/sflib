package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/sfborg/sflib/pkg/xsv"
	"github.com/stretchr/testify/assert"
)

func TestNewXsv(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewXsv()
	_, ok := sfga.(xsv.Archive)
	assert.True(ok)
}
