package coldp_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/sflib"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	res := coldp.New()
	_, ok := res.(sflib.CoLDP)
	assert.True(ok)
}
