package icoldp_test

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/internal/icoldp"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

var testDir string

func TestMain(m *testing.M) {
	setupGlobal()
	code := m.Run() // Run all tests
	teardownGlobal()
	os.Exit(code)
}

func setupGlobal() {
	var err error
	testDir, err = os.MkdirTemp("", "coldp-test")
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

func TestNameExtract(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path string
	}{
		{"name-yaml", "name/ptero-yaml.zip"},
		{"usg-yaml", "nameusage/ptero-yaml.zip"},
	}

	for _, v := range tests {
		err := gnsys.CleanDir(testDir)
		assert.Nil(err)
		path := filepath.Join("..", "..", "testdata", "coldp", v.path)
		coldp := icoldp.New()
		err = coldp.Fetch(path, testDir)
		assert.Nil(err)
		assert.NotNil(coldp)

		assert.Nil(err)
	}
}

func TestNoFile(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	coldp := icoldp.New()
	err = coldp.Fetch("notfile", testDir)
	assert.IsType(&arch.ErrFileNotFound{}, err)
}

func TestNoURL(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	coldp := icoldp.New()
	err = coldp.Fetch("http://not-a-url", testDir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrDownload{}, err)
}

func TestNotZip(t *testing.T) {
	assert := assert.New(t)
	err := gnsys.CleanDir(testDir)
	assert.Nil(err)

	path := filepath.Join("..", "..", "testdata", "notzip.zip")
	coldp := icoldp.New()
	err = coldp.Fetch(path, testDir)
	assert.NotNil(err)
	assert.IsType(&arch.ErrImportArchive{}, err)
}

func TestMeta(t *testing.T) {
	assert := assert.New(t)
	var err error
	var meta *coldp.Meta

	tests := []struct {
		msg, path, creatorGiven, contactOrganization, license string
	}{
		{"name-yaml", "name/ptero-yaml.zip", "Donald", "Species 2000", "CC0"},
		{"usg-yaml", "nameusage/ptero-yaml.zip", "Donald", "Species 2000", "cc0"},
	}

	for _, v := range tests {
		err = gnsys.CleanDir(testDir)
		assert.Nil(err)

		path := filepath.Join("..", "..", "testdata", "coldp", v.path)
		coldp := icoldp.New()
		err = coldp.Fetch(path, testDir)
		assert.Nil(err)
		err = coldp.DirInfo()
		assert.Nil(err)
		meta, err = coldp.Meta()
		assert.Nil(err)
		assert.Equal(v.creatorGiven, meta.Creators[0].Given)
		assert.Equal(v.contactOrganization, meta.Contact.Organization)
		assert.Equal(v.license, meta.License)
	}
}

func TestNoSciName(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gnsys.CleanDir(testDir)
	assert.Nil(err)

	a := icoldp.New()
	path := filepath.Join("..", "..", "testdata", "coldp", "nameusage", "no-sci-name.zip")
	err = a.Fetch(path, testDir)
	assert.Nil(err)

	err = a.DirInfo()
	assert.Nil(err)

	nuPath, ok := a.DataPaths()[coldp.NameUsageDT]
	assert.True(ok)

	ch := make(chan coldp.NameUsage)
	var res []coldp.NameUsage
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for nu := range ch {
			res = append(res, nu)
		}
	}()

	err = coldp.Read(a.Config(), nuPath, ch)
	assert.Nil(err)
	close(ch)
	wg.Wait()

	// Build lookup by ID
	byID := make(map[string]coldp.NameUsage)
	for _, nu := range res {
		byID[nu.ID] = nu
	}

	assert.Equal(11, len(res), "should load all 11 records")

	// Records with col:scientificName keep their value as-is.
	nu := byID["animalia"]
	assert.Equal("Animalia", nu.ScientificName)

	// Uninomial-only record: ScientificName assembled from col:uninomial.
	nu = byID["urn:lsid:nmbe.ch:spiderfam:0001"]
	assert.Equal("Liphistiidae", nu.ScientificName)

	nu = byID["urn:lsid:nmbe.ch:spidergen:00001"]
	assert.Equal("Heptathela", nu.ScientificName)

	// Species record: ScientificName assembled from col:genericName + col:specificEpithet.
	nu = byID["urn:lsid:nmbe.ch:spidersp:000008"]
	assert.Equal("Heptathela higoensis", nu.ScientificName)

	nu = byID["urn:lsid:nmbe.ch:spidersp:000030"]
	assert.Equal("Liphistius albipes", nu.ScientificName)
}

func TestNoSciNameSplit(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gnsys.CleanDir(testDir)
	assert.Nil(err)

	a := icoldp.New()
	path := filepath.Join("..", "..", "testdata", "coldp", "name", "no-sci-name.zip")
	err = a.Fetch(path, testDir)
	assert.Nil(err)

	err = a.DirInfo()
	assert.Nil(err)

	namePath, ok := a.DataPaths()[coldp.NameDT]
	assert.True(ok)

	taxonPath, ok := a.DataPaths()[coldp.TaxonDT]
	assert.True(ok)

	// Read Name records.
	nameCh := make(chan coldp.Name)
	var names []coldp.Name
	var nameWg sync.WaitGroup
	nameWg.Add(1)
	go func() {
		defer nameWg.Done()
		for n := range nameCh {
			names = append(names, n)
		}
	}()
	err = coldp.Read(a.Config(), namePath, nameCh)
	assert.Nil(err)
	close(nameCh)
	nameWg.Wait()

	assert.Equal(11, len(names), "should load all 11 Name records")

	byNameID := make(map[string]coldp.Name)
	for _, n := range names {
		byNameID[n.ID] = n
	}

	// Records with col:scientificName keep their value as-is.
	n := byNameID["n-animalia"]
	assert.Equal("Animalia", n.ScientificName)

	// Uninomial records: ScientificName assembled from col:uninomial.
	n = byNameID["n-liphistiidae"]
	assert.Equal("Liphistiidae", n.ScientificName)
	assert.Equal("Thorell, 1869", n.Authorship)

	n = byNameID["n-heptathela"]
	assert.Equal("Heptathela", n.ScientificName)

	// Species records: ScientificName assembled from col:genus + col:specificEpithet.
	n = byNameID["n-higoensis"]
	assert.Equal("Heptathela higoensis", n.ScientificName)
	assert.Equal("Haupt, 1983", n.Authorship)

	n = byNameID["n-albipes"]
	assert.Equal("Liphistius albipes", n.ScientificName)

	// Read Taxon records and verify nameID links.
	taxonCh := make(chan coldp.Taxon)
	var taxa []coldp.Taxon
	var taxonWg sync.WaitGroup
	taxonWg.Go(func() {
		for tx := range taxonCh {
			taxa = append(taxa, tx)
		}
	})
	err = coldp.Read(a.Config(), taxonPath, taxonCh)
	assert.Nil(err)
	close(taxonCh)
	taxonWg.Wait()

	assert.Equal(11, len(taxa), "should load all 11 Taxon records")

	byTaxonID := make(map[string]coldp.Taxon)
	for _, tx := range taxa {
		byTaxonID[tx.ID] = tx
	}

	tx := byTaxonID["animalia"]
	assert.Equal("n-animalia", tx.NameID)
	assert.Equal("", tx.ParentID)

	tx = byTaxonID["urn:lsid:nmbe.ch:spidersp:000008"]
	assert.Equal("n-higoensis", tx.NameID)
	assert.Equal("urn:lsid:nmbe.ch:spidergen:00001", tx.ParentID)
}

func TestReferences(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gnsys.CleanDir(testDir)
	assert.Nil(err)

	a := icoldp.New()
	path := filepath.Join("..", "..", "testdata", "coldp", "nameusage/ptero-yaml.zip")
	err = a.Fetch(path, testDir)
	assert.Nil(err)

	err = a.DirInfo()
	assert.Nil(err)

	refPath, ok := a.DataPaths()[coldp.ReferenceDT]
	assert.True(ok)

	ch := make(chan coldp.Reference)
	var refs []coldp.Reference
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for ref := range ch {
			refs = append(refs, ref)
		}
	}()

	err = coldp.Read(a.Config(), refPath, ch)
	assert.Nil(err)
	close(ch)
	wg.Wait()

	assert.Equal(1966, len(refs))
	assert.Equal("1597", refs[0].ID)
	assert.NotEmpty(refs[0].Citation)
}

func TestName(t *testing.T) {
	var err error
	assert := assert.New(t)
	err = gnsys.CleanDir(testDir)
	assert.Nil(err)

	a := icoldp.New()
	path := filepath.Join("..", "..", "testdata", "coldp", "name/ptero-yaml.zip")
	err = a.Fetch(path, testDir)
	assert.Nil(err)

	err = a.DirInfo()
	assert.Nil(err)

	namePath, ok := a.DataPaths()[coldp.NameDT]
	assert.True(ok)

	ch := make(chan coldp.Name)

	go func() {
		for range ch {
		}
	}()

	err = coldp.Read(a.Config(), namePath, ch)
	assert.Nil(err)
}
