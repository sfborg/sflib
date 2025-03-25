package ixsv_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/ixsv"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestXsv(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg     string
		src     string
		out     string
		errNil  bool
		rowsNum int
	}{
		{
			msg:     "csv",
			src:     "ioc-bird.csv",
			errNil:  true,
			rowsNum: 11072,
		},
		{
			msg:     "tsv",
			src:     "ioc-bird.tsv",
			errNil:  true,
			rowsNum: 9,
		},
		{
			msg:     "psv",
			src:     "ioc-bird.psv",
			errNil:  true,
			rowsNum: 9,
		},
		{
			msg:     "non-existing-file",
			src:     "non-existing.csv",
			errNil:  false,
			rowsNum: 0,
		},
	}

	dir, err := os.MkdirTemp("", "xsv-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	for _, v := range tests {
		gnsys.CleanDir(dir)
		src := filepath.Join("../../testdata/xsv", v.src)
		a := ixsv.New()

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

		err = a.Load(context.Background(), ch, 10, coldp.Cultivars)
		assert.Nil(err)
		close(ch)

		wg.Wait()
		assert.Equal(v.rowsNum, len(res), v.msg)
		assert.Equal(coldp.Zoological, res[0].Code)

	}
}
