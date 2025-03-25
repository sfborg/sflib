package itext_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/sfborg/sflib/internal/itext"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, file, name string
		errNil          bool
	}{
		{"real", "names.txt", "Plagiognathus chrysanthemi Wolff 1804", true},
		{"fake", "fake", "", false},
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

		ch := make(chan string)
		var wg sync.WaitGroup
		wg.Add(1)

		var res []string

		go func() {
			defer wg.Done()
			for v := range ch {
				res = append(res, v)
			}
		}()

		err = a.Load(context.Background(), ch)
		assert.Nil(err)
		close(ch)

		wg.Wait()
		assert.Equal(1000, len(res))
		assert.Equal(v.name, res[0])
	}
}
