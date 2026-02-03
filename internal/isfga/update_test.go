package isfga_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "update")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	path := "../../testdata/sfga/ptero_v0.3.31.sqlite"
	oldSfga := isfga.New()
	err = oldSfga.Fetch(path, dir)
	assert.Nil(err)
	_, err = oldSfga.Connect()
	assert.Nil(err)
	version := oldSfga.Version()

	newSfga := isfga.New()
	err = newSfga.Create(filepath.Join(dir))
	assert.Nil(err)

	oldSfga.Update(newSfga, false)
	_, err = newSfga.Connect()
	assert.Nil(err)
	assert.NotNil(newSfga.Db())
	var res string
	err = newSfga.Db().QueryRow("select count(*) from taxon").Scan(&res)
	assert.Nil(err)
	assert.Equal("1700", res)

	newVersion := newSfga.Version()
	// 1 means bigger
	assert.Equal(1, gnlib.CmpVersion(newVersion, version))
}

func TestUpdateWithHierarchyDiptera(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "update-diptera")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	// Use diptera-flat.sqlite which has flat hierarchy with genus/subgenus/species
	path := "../../testdata/sfga/diptera-flat.sqlite"
	oldSfga := isfga.New()
	err = oldSfga.Fetch(path, dir)
	assert.Nil(err)
	_, err = oldSfga.Connect()
	assert.Nil(err)

	// Get original taxon count
	var origCount int
	err = oldSfga.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&origCount)
	assert.Nil(err)

	newSfga := isfga.New()
	err = newSfga.Create(filepath.Join(dir))
	assert.Nil(err)

	// Update with hierarchy building enabled
	err = oldSfga.Update(newSfga, true)
	assert.Nil(err)

	_, err = newSfga.Connect()
	assert.Nil(err)

	// Check that new taxon count is greater (generated parents added)
	var newCount int
	err = newSfga.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&newCount)
	assert.Nil(err)
	assert.Greater(newCount, origCount, "should have more taxa after hierarchy build")

	// CRITICAL: Verify no self-references - a taxon should never be its own parent
	var selfRefCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__id = col__parent_id
	`).Scan(&selfRefCount)
	assert.Nil(err)
	assert.Equal(0, selfRefCount, "no taxon should reference itself as parent")

	// Verify no duplicate scientific names created
	var dupCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT n.col__scientific_name, COUNT(*) as cnt
			FROM taxon t
			JOIN name n ON t.col__name_id = n.col__id
			GROUP BY n.col__scientific_name
			HAVING cnt > 1
		)
	`).Scan(&dupCount)
	assert.Nil(err)
	assert.Equal(0, dupCount, "should have no duplicate scientific names")

	// Verify genus records don't have parent with same scientific name
	var genusWithSameParent int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__rank_id = 'GENUS'
		  AND n1.col__scientific_name = n2.col__scientific_name
	`).Scan(&genusWithSameParent)
	assert.Nil(err)
	assert.Equal(0, genusWithSameParent, "genus should not have parent with same name")

	// Verify subgenus records don't have parent with same scientific name
	var subgenusWithSameParent int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__rank_id = 'SUBGENUS'
		  AND n1.col__scientific_name = n2.col__scientific_name
	`).Scan(&subgenusWithSameParent)
	assert.Nil(err)
	assert.Equal(0, subgenusWithSameParent, "subgenus should not have parent with same name")

	// Verify hierarchy integrity: all parent_ids should exist as taxon ids
	var orphanCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		WHERE t1.col__parent_id != ''
		  AND t1.col__parent_id NOT IN (SELECT col__id FROM taxon)
	`).Scan(&orphanCount)
	assert.Nil(err)
	assert.Equal(0, orphanCount, "all parent_ids should reference existing taxa")

	// Verify a specific genus (Amphineurus) has Subfamily as parent (closest in hierarchy)
	var amphineurusParentRank string
	err = newSfga.Db().QueryRow(`
		SELECT n2.col__rank_id FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__scientific_name = 'Amphineurus'
		  AND n1.col__rank_id = 'GENUS'
	`).Scan(&amphineurusParentRank)
	assert.Nil(err)
	assert.Equal("SUBFAMILY", amphineurusParentRank, "genus Amphineurus should have subfamily as parent")

	// Verify species with subgenus has subgenus as parent, not genus
	// "Amphineurus (Amphineurus) breviclavus" should have parent "Amphineurus (Amphineurus)" (SUBGENUS)
	var speciesParentName, speciesParentRank string
	err = newSfga.Db().QueryRow(`
		SELECT n2.col__scientific_name, n2.col__rank_id FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__scientific_name = 'Amphineurus (Amphineurus) breviclavus'
	`).Scan(&speciesParentName, &speciesParentRank)
	assert.Nil(err)
	assert.Equal("Amphineurus (Amphineurus)", speciesParentName, "species should have subgenus as parent")
	assert.Equal("SUBGENUS", speciesParentRank, "parent should be SUBGENUS rank")

	// Verify subgenus has genus as parent
	var subgenusParentName, subgenusParentRank string
	err = newSfga.Db().QueryRow(`
		SELECT n2.col__scientific_name, n2.col__rank_id FROM taxon t1
		JOIN name n1 ON t1.col__name_id = n1.col__id
		JOIN taxon t2 ON t1.col__parent_id = t2.col__id
		JOIN name n2 ON t2.col__name_id = n2.col__id
		WHERE n1.col__scientific_name = 'Amphineurus (Amphineurus)'
		  AND n1.col__rank_id = 'SUBGENUS'
	`).Scan(&subgenusParentName, &subgenusParentRank)
	assert.Nil(err)
	assert.Equal("Amphineurus", subgenusParentName, "subgenus should have genus as parent")
	assert.Equal("GENUS", subgenusParentRank, "parent should be GENUS rank")

	// Verify SubgenusID points to the subgenus record, not the genus
	var subgenusID, genusID string
	err = newSfga.Db().QueryRow(`
		SELECT t.sf__subgenus_id, t.sf__genus_id FROM taxon t
		JOIN name n ON t.col__name_id = n.col__id
		WHERE n.col__scientific_name = 'Amphineurus (Amphineurus) breviclavus'
	`).Scan(&subgenusID, &genusID)
	assert.Nil(err)
	assert.NotEqual(subgenusID, genusID, "SubgenusID and GenusID should differ")
	// SubgenusID should match the subgenus record's ID
	var actualSubgenusRecordID string
	err = newSfga.Db().QueryRow(`
		SELECT t.col__id FROM taxon t
		JOIN name n ON t.col__name_id = n.col__id
		WHERE n.col__scientific_name = 'Amphineurus (Amphineurus)'
		  AND n.col__rank_id = 'SUBGENUS'
	`).Scan(&actualSubgenusRecordID)
	assert.Nil(err)
	assert.Equal(actualSubgenusRecordID, subgenusID, "SubgenusID should point to the subgenus record")
}

func TestUpdateWithHierarchy(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join(testDir, "update-hier")
	err := os.Mkdir(dir, 0755)
	assert.Nil(err)
	defer os.RemoveAll(dir)

	// Use virus-flat-hier.sqlite which has flat hierarchy but no parent/child
	path := "../../testdata/sfga/virus-flat-hier.sqlite"
	oldSfga := isfga.New()
	err = oldSfga.Fetch(path, dir)
	assert.Nil(err)
	_, err = oldSfga.Connect()
	assert.Nil(err)

	// Verify source has no parent_id set
	var noParents int
	err = oldSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__parent_id IS NOT NULL AND col__parent_id != ''
	`).Scan(&noParents)
	assert.Nil(err)
	assert.Equal(0, noParents, "source should have no parent_id set")

	// Get original taxon count
	var origCount int
	err = oldSfga.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&origCount)
	assert.Nil(err)

	newSfga := isfga.New()
	err = newSfga.Create(filepath.Join(dir))
	assert.Nil(err)

	// Update with hierarchy building enabled
	err = oldSfga.Update(newSfga, true)
	assert.Nil(err)

	_, err = newSfga.Connect()
	assert.Nil(err)

	// Check that new taxon count is greater (generated parents added)
	var newCount int
	err = newSfga.Db().QueryRow("SELECT COUNT(*) FROM taxon").Scan(&newCount)
	assert.Nil(err)
	assert.Greater(newCount, origCount, "should have more taxa after hierarchy build")

	// Check that parent_id is now set for original taxa
	var hasParents int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__parent_id IS NOT NULL AND col__parent_id != ''
	`).Scan(&hasParents)
	assert.Nil(err)
	assert.Greater(hasParents, 0, "should have parent_id set after hierarchy build")

	// Check that generated taxa have sf-{int} format (not sf-{uuid})
	var sfCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon WHERE col__id LIKE 'sf-%'
	`).Scan(&sfCount)
	assert.Nil(err)
	assert.Greater(sfCount, 0, "should have generated sf- prefixed taxa")

	// Verify IDs are sequential integers, not UUIDs
	var sampleID string
	err = newSfga.Db().QueryRow(`
		SELECT col__id FROM taxon WHERE col__id LIKE 'sf-%' LIMIT 1
	`).Scan(&sampleID)
	assert.Nil(err)
	// sf-{int} format should be short (e.g., "sf-1", "sf-42"), not sf-{uuid} which is 39 chars
	assert.Less(len(sampleID), 15, "ID should be sf-{int} format, not sf-{uuid}")

	// Verify generated parent taxa have their classification IDs set
	var familyWithIDs int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE col__id LIKE 'sf-%'
		  AND sf__family_id != ''
		  AND sf__kingdom_id != ''
	`).Scan(&familyWithIDs)
	assert.Nil(err)
	assert.Greater(familyWithIDs, 0, "generated parent taxa should have classification IDs set")

	// Verify Realm is populated (virus data has Realm)
	var realmCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon
		WHERE sf__realm != '' AND sf__realm_id != ''
	`).Scan(&realmCount)
	assert.Nil(err)
	assert.Greater(realmCount, 0, "taxa should have Realm and RealmID set")

	// Verify hierarchy integrity: all parent_ids should exist as taxon ids
	var orphanCount int
	err = newSfga.Db().QueryRow(`
		SELECT COUNT(*) FROM taxon t1
		WHERE t1.col__parent_id != ''
		  AND t1.col__parent_id NOT IN (SELECT col__id FROM taxon)
	`).Scan(&orphanCount)
	assert.Nil(err)
	assert.Equal(0, orphanCount, "all parent_ids should reference existing taxa")
}
