package itext_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/sfborg/sflib/internal/itext"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, file string
		errNil    bool
	}{
		{"real", "names.txt", true},
		{"fake", "fake", false},
	}

	testDir, err := os.MkdirTemp("", "sflib-text")
	assert.Nil(err)
	defer os.RemoveAll(testDir)

	for _, v := range tests {
		a := itext.New()
		path := filepath.Join("../../testdata/text", v.file)
		err := a.Fetch(path, testDir)
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
		assert.Equal(1000, len(res))
		assert.Equal(coldp.Cultivars, res[0].Code)
	}
}
