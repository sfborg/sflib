package isfga

import (
	"context"
	"database/sql"
	"testing"

	atlasschema "ariga.io/atlas/sql/schema"
	atlassqlite "ariga.io/atlas/sql/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestAtlasPureGo verifies that Atlas works with modernc.org/sqlite (pure Go,
// no CGO) and can inspect schemas, compute diffs, and apply changes.
func TestAtlasPureGo(t *testing.T) {
	ctx := context.Background()

	// Create an "old" database with a simple schema.
	oldDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer oldDB.Close()

	_, err = oldDB.Exec(`
		CREATE TABLE names (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		);
		CREATE TABLE version (
			id INTEGER PRIMARY KEY,
			version TEXT NOT NULL
		);
	`)
	require.NoError(t, err)

	// Create a "desired" database with an extra column and table.
	desiredDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer desiredDB.Close()

	_, err = desiredDB.Exec(`
		CREATE TABLE names (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			authorship TEXT DEFAULT ''
		);
		CREATE TABLE version (
			id INTEGER PRIMARY KEY,
			version TEXT NOT NULL
		);
		CREATE TABLE refs (
			id TEXT PRIMARY KEY,
			citation TEXT
		);
	`)
	require.NoError(t, err)

	// Open Atlas drivers on both.
	oldDrv, err := atlassqlite.Open(oldDB)
	require.NoError(t, err)

	desiredDrv, err := atlassqlite.Open(desiredDB)
	require.NoError(t, err)

	// Inspect both schemas.
	current, err := oldDrv.InspectSchema(ctx, "main", nil)
	require.NoError(t, err)

	desired, err := desiredDrv.InspectSchema(ctx, "main", nil)
	require.NoError(t, err)

	// Compute diff.
	changes, err := oldDrv.SchemaDiff(current, desired)
	require.NoError(t, err)
	assert.True(t, len(changes) > 0, "expected schema changes")

	// Apply changes to the old DB.
	err = oldDrv.ApplyChanges(ctx, changes)
	require.NoError(t, err)

	// Verify: inspect old DB again — should now match desired.
	updated, err := oldDrv.InspectSchema(ctx, "main", nil)
	require.NoError(t, err)

	noChanges, err := oldDrv.SchemaDiff(updated, desired)
	require.NoError(t, err)
	assert.Empty(t, noChanges, "schemas should match after migration")

	// Verify the new column exists.
	var authorship string
	err = oldDB.QueryRow(`
		INSERT INTO names (id, name, authorship) VALUES ('1', 'Homo sapiens', 'Linnaeus')
		RETURNING authorship
	`).Scan(&authorship)
	require.NoError(t, err)
	assert.Equal(t, "Linnaeus", authorship)

	// Verify the new table exists.
	_, err = oldDB.Exec(`INSERT INTO refs (id, citation) VALUES ('1', 'Some ref')`)
	require.NoError(t, err)

	// Log the changes that were applied.
	for _, c := range changes {
		switch c := c.(type) {
		case *atlasschema.AddTable:
			t.Logf("Added table: %s", c.T.Name)
		case *atlasschema.ModifyTable:
			t.Logf("Modified table: %s", c.T.Name)
		case *atlasschema.DropTable:
			t.Logf("Dropped table: %s", c.T.Name)
		}
	}
}
