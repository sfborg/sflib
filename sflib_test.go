package sflib_test

import (
	"testing"

	"github.com/sfborg/sflib"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/sfga"
	"github.com/sfborg/sflib/pkg/text"
	"github.com/sfborg/sflib/pkg/xsv"
	"github.com/stretchr/testify/assert"
)

func TestNewText(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewText()
	_, ok := sfga.(text.Archive)
	assert.True(ok)
}

func TestNewXsv(t *testing.T) {
	assert := assert.New(t)
	sfga := sflib.NewXsv()
	_, ok := sfga.(xsv.Archive)
	assert.True(ok)
}

func TestNewColdp(t *testing.T) {
	assert := assert.New(t)
	res := sflib.NewColdp()
	_, ok := res.(coldp.Archive)
	assert.True(ok)
}

func TestNewSfga(t *testing.T) {
	assert := assert.New(t)
	res := sflib.NewSfga()
	_, ok := res.(sfga.Archive)
	assert.True(ok)
}
