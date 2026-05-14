package isfga_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "migrate")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/ptero_v0.5.1.sqlite"
	oldSfga := isfga.New()
	err = oldSfga.Fetch(path, dir)
	assert.Nil(err)

	// After Fetch, auto-migration should have run.
	// The version should now be current.
	newVersion := oldSfga.Version()
	// 0 means equal, 1 means file is newer — either is acceptable.
	assert.NotEqual(-1, gnlib.CmpVersion(newVersion, "v0.3.31"),
		"migrated version should be >= v0.3.31")

	// Taxon count should be preserved.
	_, err = oldSfga.Connect()
	assert.Nil(err)
	var count string
	err = oldSfga.Db().QueryRow("SELECT count(*) FROM taxon").Scan(&count)
	assert.Nil(err)
	assert.Equal("1700", count)
}

func TestMigrateExplicit(t *testing.T) {
	require := require.New(t)
	dir := filepath.Join(testDir, "migrate-explicit")
	err := os.Mkdir(dir, 0755)
	require.NoError(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/ptero_v0.5.1.sqlite"
	src := isfga.New()
	err = src.Fetch(path, dir)
	require.NoError(err)

	migDir := filepath.Join(dir, "mig-out")
	migrated, err := src.Migrate(migDir)
	require.NoError(err)
	require.NotNil(migrated)

	// Migrated archive should be at current schema version.
	srcVersion := src.Version()
	migVersion := migrated.Version()
	require.GreaterOrEqual(gnlib.CmpVersion(migVersion, srcVersion), 0,
		"migrated version should be >= source version")

	// Data should be intact.
	_, err = migrated.Connect()
	require.NoError(err)
	var count int
	err = migrated.Db().QueryRow("SELECT count(*) FROM taxon").Scan(&count)
	require.NoError(err)
	require.Equal(1700, count)
}

func TestAddParentsDiptera(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "add-parents-diptera")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/diptera-flat.sqlite"
	a := isfga.New()
	err = a.Fetch(path, dir)
	assert.Nil(err)
	_, err = a.Connect()
	assert.Nil(err)

	var origCount int
	err = a.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&origCount)
	assert.Nil(err)

	err = a.AddParents(context.Background())
	assert.Nil(err)

	var newCount int
	err = a.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&newCount)
	assert.Nil(err)
	assert.Greater(newCount, origCount, "should have more taxa after hierarchy build")

	// No self-references.
	var selfRefCount int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon WHERE col__id = col__parent_id
	`).Scan(&selfRefCount)
	assert.Nil(err)
	assert.Equal(0, selfRefCount, "no taxon should reference itself as parent")

	// No duplicate scientific names.
	var dupCount int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT n.col__scientific_name, COUNT(*) as cnt
			FROM taxon t JOIN name n ON t.col__name_id = n.col__id
			GROUP BY n.col__scientific_name HAVING cnt > 1
		)
	`).Scan(&dupCount)
	assert.Nil(err)
	assert.Equal(0, dupCount, "no duplicate scientific names")

	// All parent_ids reference existing taxa.
	var orphanCount int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		WHERE t1.col__parent_id != ''
		  AND t1.col__parent_id NOT IN (SELECT col__id FROM taxon)
	`).Scan(&orphanCount)
	assert.Nil(err)
	assert.Equal(0, orphanCount, "all parent_ids should reference existing taxa")

	// Genus Amphineurus should have Subfamily as parent.
	var amphineurusParentRank string
	err = a.Db().QueryRow(`
		SELECT n2.col__rank_id FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__scientific_name = 'Amphineurus' AND n1.col__rank_id = 'GENUS'
	`).Scan(&amphineurusParentRank)
	assert.Nil(err)
	assert.Equal("SUBFAMILY", amphineurusParentRank)
}

func TestAddParentsVirus(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "add-parents-virus")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/virus-flat-hier.sqlite"
	a := isfga.New()
	err = a.Fetch(path, dir)
	assert.Nil(err)
	_, err = a.Connect()
	assert.Nil(err)

	var noParents int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__parent_id IS NOT NULL AND col__parent_id != ''
	`).Scan(&noParents)
	assert.Nil(err)
	assert.Equal(0, noParents, "source should have no parent_id set")

	var origCount int
	err = a.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&origCount)
	assert.Nil(err)

	err = a.AddParents(context.Background())
	assert.Nil(err)

	var newCount int
	err = a.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&newCount)
	assert.Nil(err)
	assert.Greater(newCount, origCount, "should have more taxa after hierarchy build")

	var hasParents int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__parent_id IS NOT NULL AND col__parent_id != ''
	`).Scan(&hasParents)
	assert.Nil(err)
	assert.Greater(hasParents, 0, "should have parent_id set after hierarchy build")

	// Hierarchy integrity.
	var orphanCount int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		WHERE t1.col__parent_id != ''
		  AND t1.col__parent_id NOT IN (SELECT col__id FROM taxon)
	`).Scan(&orphanCount)
	assert.Nil(err)
	assert.Equal(0, orphanCount, "all parent_ids should reference existing taxa")

	// Realm should be populated.
	var realmCount int
	err = a.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon WHERE sf__realm != '' AND sf__realm_id != ''
	`).Scan(&realmCount)
	assert.Nil(err)
	assert.Greater(realmCount, 0, "taxa should have Realm and RealmID set")
}
