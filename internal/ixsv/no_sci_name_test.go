package ixsv_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/internal/ixsv"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestNoSciName(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, src  string
		parseCode nomcode.Code
		errNil    bool
		rowsNum   int
	}{
		{
			msg:       "virus csv",
			parseCode: nomcode.Virus,
			src:       "viruses-no-sci-name.csv",
			errNil:    true,
			rowsNum:   84,
		},
	}

	dir, err := os.MkdirTemp("", "xsv-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	for _, v := range tests {
		gnsys.CleanDir(dir)
		src := filepath.Join("../../testdata/xsv", v.src)
		a := ixsv.New(config.OptNomCode(v.parseCode))

		err = a.Fetch(src, dir)
		assert.Equal(v.errNil, err == nil)

		if err != nil {
			continue
		}

		ch := make(chan coldp.NameUsage)
		var wg sync.WaitGroup
		wg.Add(1)

		var res []coldp.NameUsage

		go func() {
			defer wg.Done()
			for v := range ch {
				res = append(res, v)
			}
		}()

		err = a.Load(context.Background(), ch, 10, v.parseCode)
		assert.Nil(err)
		close(ch)

		wg.Wait()
		assert.Equal(v.rowsNum, len(res), v.msg)
		assert.Equal(v.parseCode, res[0].Code)
	}
}
