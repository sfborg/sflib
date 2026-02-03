package ixsv_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/isfga"
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
		code    nomcode.Code
	}{
		{
			msg:     "csv",
			src:     "ioc-bird.csv",
			errNil:  true,
			rowsNum: 11072,
			code:    nomcode.Zoological,
		},
		{
			msg:     "tsv",
			src:     "ioc-bird.tsv",
			errNil:  true,
			rowsNum: 9,
			code:    nomcode.Zoological,
		},
		{
			msg:     "psv",
			src:     "ioc-bird.psv",
			errNil:  true,
			rowsNum: 9,
			code:    nomcode.Zoological,
		},
		{
			msg:     "diptera",
			src:     "diptera.csv",
			errNil:  true,
			rowsNum: 16274,
			code:    nomcode.Zoological,
		},
		{
			msg:     "viruses",
			src:     "viruses-no-sci-name.csv",
			errNil:  true,
			rowsNum: 84,
			code:    nomcode.Virus,
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

		err = a.Load(context.Background(), ch, 10, v.code)
		assert.Nil(err)
		close(ch)

		wg.Wait()
		assert.Equal(v.rowsNum, len(res), v.msg)
		assert.Equal(v.code, res[0].Code)

	}
}

func TestDipteraClassification(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "xsv-diptera-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	src := filepath.Join("../../testdata/xsv", "diptera.csv")
	a := ixsv.New()

	err = a.Fetch(src, dir)
	assert.Nil(err)

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

	err = a.Load(context.Background(), ch, 1, nomcode.Zoological)
	assert.Nil(err)
	close(ch)

	wg.Wait()

	// Build a map by scientificName for easier lookup
	byName := make(map[string]coldp.NameUsage)
	for _, nu := range res {
		byName[nu.ScientificName] = nu
	}

	// Test family rank record
	nu := byName["Limoniidae"]
	assert.Equal("Limoniidae", nu.Family, "Family should be set")
	assert.Equal(coldp.Family, nu.Rank, "Rank should be family")

	// Test genus rank record
	nu = byName["Amphineurus"]
	assert.Equal("Limoniidae", nu.Family, "Family should be Limoniidae")
	assert.Equal("Chioneinae", nu.Subfamily, "Subfamily should be Chioneinae")
	assert.Equal("Amphineurus", nu.Genus, "Genus should be Amphineurus")
	assert.Equal(coldp.Genus, nu.Rank, "Rank should be genus")

	// Test subgenus rank record
	nu = byName["Amphineurus (Amphineurus)"]
	assert.Equal("Limoniidae", nu.Family, "Family should be Limoniidae")
	assert.Equal("Chioneinae", nu.Subfamily, "Subfamily should be Chioneinae")
	assert.Equal("Amphineurus", nu.Genus, "Genus should be Amphineurus")
	assert.Equal("Amphineurus", nu.Subgenus, "Subgenus should be Amphineurus")
	assert.Equal(coldp.Subgenus, nu.Rank, "Rank should be subgenus")

	// Test species rank record
	nu = byName["Amphineurus (Amphineurus) bicinctus"]
	assert.Equal("Limoniidae", nu.Family, "Family should be Limoniidae")
	assert.Equal("Chioneinae", nu.Subfamily, "Subfamily should be Chioneinae")
	assert.Equal("Amphineurus", nu.Genus, "Genus should be Amphineurus")
	assert.Equal("Amphineurus", nu.Subgenus, "Subgenus should be Amphineurus")
	assert.Equal("bicinctus", nu.Species, "Species should be bicinctus")
	assert.Equal("Edwards, 1923", nu.Authorship, "Authorship should be set")
	assert.Equal(coldp.Species, nu.Rank, "Rank should be species")
}

func TestVirusClassification(t *testing.T) {
	assert := assert.New(t)

	dir, err := os.MkdirTemp("", "xsv-virus-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	src := filepath.Join("../../testdata/xsv", "viruses-no-sci-name.csv")
	a := ixsv.New()

	err = a.Fetch(src, dir)
	assert.Nil(err)

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

	err = a.Load(context.Background(), ch, 1, nomcode.Virus)
	assert.Nil(err)
	close(ch)

	wg.Wait()

	// Build a map by scientificName for easier lookup
	byName := make(map[string]coldp.NameUsage)
	for _, nu := range res {
		byName[nu.ScientificName] = nu
	}

	// Test virus with full classification hierarchy
	// Row: 1,Adnaviria,,Zilligvirae,,Taleaviricota,,Tokiviricetes,,Ligamenvirales,,Chiyouviridae,,Wargodvirus,,Wargodvirus xiongnu
	nu := byName["Wargodvirus xiongnu"]
	assert.Equal("Zilligvirae", nu.Kingdom, "Kingdom should be Zilligvirae")
	assert.Equal("Taleaviricota", nu.Phylum, "Phylum should be Taleaviricota")
	assert.Equal("Tokiviricetes", nu.Class, "Class should be Tokiviricetes")
	assert.Equal("Ligamenvirales", nu.Order, "Order should be Ligamenvirales")
	assert.Equal("Chiyouviridae", nu.Family, "Family should be Chiyouviridae")
	assert.Equal("Wargodvirus", nu.Genus, "Genus should be Wargodvirus")
	assert.Equal(nomcode.Virus, nu.Code, "Code should be Virus")

	// Test virus with subfamily
	// Row: 177,Duplodnaviria,,Heunggongvirae,,Uroviricota,,Caudoviricetes,,Autographivirales,,Autonotataviridae,Gujervirinae,Pradovirus,,Pradovirus F5
	nu = byName["Pradovirus F5"]
	assert.Equal("Heunggongvirae", nu.Kingdom, "Kingdom should be Heunggongvirae")
	assert.Equal("Uroviricota", nu.Phylum, "Phylum should be Uroviricota")
	assert.Equal("Caudoviricetes", nu.Class, "Class should be Caudoviricetes")
	assert.Equal("Autographivirales", nu.Order, "Order should be Autographivirales")
	assert.Equal("Autonotataviridae", nu.Family, "Family should be Autonotataviridae")
	assert.Equal("Gujervirinae", nu.Subfamily, "Subfamily should be Gujervirinae")
	assert.Equal("Pradovirus", nu.Genus, "Genus should be Pradovirus")

	// Test virus with minimal classification (no order)
	// Row: 4692,Duplodnaviria,,Heunggongvirae,,Uroviricota,,Caudoviricetes,,,,,,Casadabanvirus,,Casadabanvirus R24
	nu = byName["Casadabanvirus R24"]
	assert.Equal("Heunggongvirae", nu.Kingdom, "Kingdom should be Heunggongvirae")
	assert.Equal("Uroviricota", nu.Phylum, "Phylum should be Uroviricota")
	assert.Equal("Caudoviricetes", nu.Class, "Class should be Caudoviricetes")
	assert.Equal("", nu.Order, "Order should be empty")
	assert.Equal("", nu.Family, "Family should be empty")
	assert.Equal("Casadabanvirus", nu.Genus, "Genus should be Casadabanvirus")
}

func TestWrite(t *testing.T) {
	assert := assert.New(t)

	// Create temp directory for output
	dir, err := os.MkdirTemp("", "xsv-write-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	// Setup SFGA archive from dinof.sqlite
	sfgaDir := filepath.Join(dir, "sfga")
	err = os.Mkdir(sfgaDir, 0755)
	assert.Nil(err)

	sfgaSrc := filepath.Join("../../testdata/sfga", "virus-flat-hier.sqlite")
	sfga := isfga.New()
	err = sfga.Fetch(sfgaSrc, sfgaDir)
	assert.Nil(err)

	_, err = sfga.Connect()
	assert.Nil(err)
	defer sfga.Close()

	// Create channel and start reading NameUsages from SFGA
	ch := make(chan coldp.NameUsage)
	var wg sync.WaitGroup

	// Start SFGA reader in goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		err := sfga.LoadNameUsages(context.Background(), ch)
		assert.Nil(err)
	}()

	// Write to CSV using ixsv
	csvPath := filepath.Join(dir, "output.csv")
	a := ixsv.New()
	err = a.Write(context.Background(), ch, csvPath)
	assert.Nil(err)

	wg.Wait()

	// Verify the output file exists
	assert.Equal(csvPath, a.FilePath())

	// Read and verify output
	f, err := os.Open(csvPath)
	assert.Nil(err)
	defer f.Close()

	// Count lines (header + data rows)
	content, err := os.ReadFile(csvPath)
	assert.Nil(err)

	lines := strings.Split(string(content), "\n")
	// Remove empty last line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Should have header + some data rows (no synonyms)
	assert.Greater(len(lines), 1, "Should have header and at least one data row")

	// Verify header starts with expected columns
	assert.True(strings.HasPrefix(lines[0], "col:id,"), "Header should start with col:id")
}

func TestWriteSkipsSynonyms(t *testing.T) {
	assert := assert.New(t)

	// Create temp directory for output
	dir, err := os.MkdirTemp("", "xsv-write-syn-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	// Setup SFGA archive from ptero database (has synonyms)
	sfgaDir := filepath.Join(dir, "sfga")
	err = os.Mkdir(sfgaDir, 0755)
	assert.Nil(err)

	sfgaSrc := filepath.Join("../../testdata/sfga", "ptero_v0.4.1.sqlite.sqlite")
	sfga := isfga.New()
	err = sfga.Fetch(sfgaSrc, sfgaDir)
	assert.Nil(err)

	_, err = sfga.Connect()
	assert.Nil(err)
	defer sfga.Close()

	// Create channel and start reading NameUsages from SFGA
	ch := make(chan coldp.NameUsage)
	var wg sync.WaitGroup

	// Start SFGA reader in goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		err := sfga.LoadNameUsages(context.Background(), ch)
		assert.Nil(err)
	}()

	// Write to CSV using ixsv
	csvPath := filepath.Join(dir, "output.csv")
	a := ixsv.New()
	err = a.Write(context.Background(), ch, csvPath)
	assert.Nil(err)

	wg.Wait()

	// Read and count lines
	content, err := os.ReadFile(csvPath)
	assert.Nil(err)

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// ptero has 1700 taxa and 1282 synonyms
	// Output should have header + 1700 data rows (synonyms skipped)
	dataRows := len(lines) - 1 // subtract header
	assert.Equal(1700, dataRows, "Should have 1700 rows (taxa only, no synonyms)")
}

func TestExport(t *testing.T) {
	assert := assert.New(t)

	// Create temp directory for output
	dir, err := os.MkdirTemp("", "xsv-export-")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	// Setup SFGA archive
	sfgaDir := filepath.Join(dir, "sfga")
	err = os.Mkdir(sfgaDir, 0755)
	assert.Nil(err)

	sfgaSrc := filepath.Join("../../testdata/sfga", "virus-flat-hier.sqlite")
	sfga := isfga.New()
	err = sfga.Fetch(sfgaSrc, sfgaDir)
	assert.Nil(err)

	_, err = sfga.Connect()
	assert.Nil(err)
	defer sfga.Close()

	// Write CSV
	ch := make(chan coldp.NameUsage)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(ch)
		sfga.LoadNameUsages(context.Background(), ch)
	}()

	csvPath := filepath.Join(dir, "data.csv")
	a := ixsv.New()
	err = a.Write(context.Background(), ch, csvPath)
	assert.Nil(err)
	wg.Wait()

	// Test Export without zip
	exportPath := filepath.Join(dir, "exported.csv")
	err = a.Export(exportPath, false)
	assert.Nil(err)

	exists, _ := gnsys.FileExists(exportPath)
	assert.True(exists, "Exported CSV should exist")

	// Test Export with zip
	exportZipPath := filepath.Join(dir, "exported_zip.csv")
	err = a.Export(exportZipPath, true)
	assert.Nil(err)

	exists, _ = gnsys.FileExists(exportZipPath)
	assert.True(exists, "Exported CSV should exist")

	exists, _ = gnsys.FileExists(exportZipPath + ".zip")
	assert.True(exists, "Exported ZIP should exist")
}
