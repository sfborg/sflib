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

func TestInatLoad(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "sflib-inat-load")

	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := filepath.Join("../../testdata/dwca", "inat-test.zip")
	a := sflib.NewDwca()

	err = a.Fetch(path, dir)
	assert.Nil(err)

	ch := make(chan coldp.Data)
	var wg sync.WaitGroup
	wg.Add(1)

	var coreCount int
	go func() {
		defer wg.Done()
		for d := range ch {
			coreCount += len(d.NameUsages)
		}
	}()

	err = a.LoadCore(context.Background(), ch)
	assert.Nil(err)
	close(ch)
	wg.Wait()
	assert.Equal(10000, coreCount)

	chVern := make(chan []coldp.Vernacular)
	wg.Add(1)

	var vernCount int
	go func() {
		defer wg.Done()
		for v := range chVern {
			vernCount += len(v)
		}
	}()

	err = a.LoadVernacular(context.Background(), 0, chVern)
	assert.Nil(err)
	close(chVern)
	wg.Wait()
	assert.Equal(12209, vernCount)
}

func TestInatToSfga(t *testing.T) {
	assert := assert.New(t)
	dwcaDir := filepath.Join(testDir, "sflib-inat-to-sfga-dwca")
	sfgaDir := filepath.Join(testDir, "sflib-inat-to-sfga-sfga")

	err := os.Mkdir(dwcaDir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dwcaDir)

	err = os.Mkdir(sfgaDir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(sfgaDir)

	dwcaArch := sflib.NewDwca()
	err = dwcaArch.Fetch(filepath.Join("../../testdata/dwca", "inat-test.zip"), dwcaDir)
	assert.Nil(err)

	sfgaArch := sflib.NewSfga()
	err = sfgaArch.Create(sfgaDir)
	assert.Nil(err)
	_, err = sfgaArch.Connect()
	assert.Nil(err)
	defer sfgaArch.Close()

	ch := make(chan coldp.Data)
	var wg sync.WaitGroup
	var insertErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		for d := range ch {
			if e := sfgaArch.InsertNameUsages(d.NameUsages); e != nil {
				insertErr = e
			}
		}
	}()

	err = dwcaArch.LoadCore(context.Background(), ch)
	assert.Nil(err)
	close(ch)
	wg.Wait()
	assert.Nil(insertErr)

	chVern := make(chan []coldp.Vernacular)
	var insertVernErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range chVern {
			if e := sfgaArch.InsertVernaculars(v); e != nil {
				insertVernErr = e
			}
		}
	}()

	err = dwcaArch.LoadVernacular(context.Background(), 0, chVern)
	assert.Nil(err)
	close(chVern)
	wg.Wait()
	assert.Nil(insertVernErr)

	db := sfgaArch.Db()
	assert.NotNil(db)

	var count int
	err = db.QueryRow("select count(*) from name").Scan(&count)
	assert.Nil(err)
	assert.Equal(10000, count)

	err = db.QueryRow("select count(*) from taxon").Scan(&count)
	assert.Nil(err)
	assert.Equal(9998, count)

	err = db.QueryRow("select count(*) from synonym").Scan(&count)
	assert.Nil(err)
	assert.Equal(0, count)

	err = db.QueryRow("select count(*) from vernacular").Scan(&count)
	assert.Nil(err)
	assert.Equal(12209, count)
}

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
