package idwca_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/sfborg/sflib"
	"github.com/sfborg/sflib/config"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setupGlobal()
	code := m.Run() // Run all tests
	teardownGlobal()
	os.Exit(code)
}

func setupGlobal() {
	var err error
	testDir, err = os.MkdirTemp("", "dwca-test")
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func teardownGlobal() {
	var err error
	err = os.RemoveAll(testDir)
	if err != nil {
		panic(err)
	}
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "sflib-dwca")

	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	tests := []struct {
		msg, file          string
		taxons, vern, dist int
	}{
		{"col", "col-mini.zip", 79, 19, 18},
		// two vernacular files
		{"2vern", "two-vern.tar.gz", 79, 19, 18},
	}

	for _, v := range tests {
		path := filepath.Join("../../testdata/dwca", v.file)
		opts := []config.Option{
			config.OptBatchSize(5),
		}
		a := sflib.NewDwca(opts...)

		err := a.Fetch(path, dir)
		assert.Nil(err)

		ch := make(chan coldp.Data)
		var wg sync.WaitGroup
		wg.Add(1)

		var coreCount int

		go func() {
			defer wg.Done()
			for d := range ch {
				recs := len(d.NameUsages)
				coreCount += recs
			}
		}()

		err = a.LoadCore(context.Background(), ch)
		assert.Nil(err)
		close(ch)

		wg.Wait()
		assert.Equal(v.taxons, coreCount)

		chVern := make(chan []coldp.Vernacular)
		wg.Add(1)

		var vernCount int
		go func() {
			defer wg.Done()
			for v := range chVern {
				recs := len(v)
				vernCount += recs
			}
		}()

		err = a.LoadVernacular(context.Background(), 2, chVern)
		assert.Nil(err)
		close(chVern)
		wg.Wait()
		assert.Equal(v.vern, vernCount)

		chDist := make(chan coldp.Data)
		wg.Add(1)

		var distCount int
		go func() {
			defer wg.Done()
			for d := range chDist {
				recs := len(d.Distributions)
				distCount += recs
			}
		}()

		err = a.LoadDistribution(context.Background(), 0, chDist)
		assert.Nil(err)
		close(chDist)
		wg.Wait()
		assert.Equal(v.dist, distCount)
	}
}
