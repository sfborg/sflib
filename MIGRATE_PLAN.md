# SFGA Schema Migration Plan

## Problem

Every SFGA file carries a schema version in its `version` table (readable via
`Archive.Version()`). The schema itself is fetched from a pinned git tag of the
`sfborg/sfga` repository at build time (see `config/config.go`: `repoTag`,
`schemaHash`).

When a new version of `sflib` is released with a newer schema tag, any SFGA
file produced by an older version becomes **outdated**. Consumers of that file
— `sf diff`, `sf apply`, `sf to`, and anything that calls `sfga.Archive`
readers — receive wrong results or silent failures because the SQL queries
reference columns or tables that may not exist in the old schema, or may
disappear in the newer one.

The problem is not in `sf from` subcommands: they always call `Archive.Create()`
which builds a fresh database at the current schema version. The problem is
exclusively at **consumption points** where an existing SFGA file is opened and
read.

## Current Migration Mechanism (`isfga.Update`)

`internal/isfga/update.go` implements `Update(emptySfga Archive, withParents
bool) error`. The approach:

1. Creates a brand-new empty database at the current schema version.
2. Introspects both databases via `PRAGMA table_info`.
3. Copies rows table-by-table using only the column names common to both
   schemas (field intersection).
4. The caller discards the old database and uses the new one.

This is exposed in `sf` only through the explicit `sf update` command
(`cmd/update.go`), which requires two arguments: old path and new path. There
is no automatic migration at consumption points.

### Limitations of the current approach

- **Data duplication**: the entire dataset is physically copied to a new file.
  For large archives this is slow and doubles disk usage temporarily.
- **No auto-migration**: `sf diff`, `sf apply`, `sf to` do not call
  `IsCompatible()` before querying. They proceed silently on stale schemas and
  may return corrupt results.
- **Manual introspection**: the field-intersection logic re-implements what a
  proper migration library would do automatically, and it cannot handle column
  renames, type changes, or removed columns.
- **`withParents` coupling**: the data-copy path is entangled with the
  hierarchy-building feature, making the code harder to reason about.

## Proposed Solution: Atlas Declarative Migration

### Why Atlas

The `ariga.io/atlas` library has a **SQLite driver**
(`ariga.io/atlas/sql/sqlite`) that provides:

- `Driver.InspectSchema()` — reads the live schema of a SQLite file into an
  in-memory `schema.Schema` struct.
- `Driver.SchemaDiff(current, desired)` — computes the minimal set of changes
  needed to bring `current` up to `desired`.
- `Driver.ApplyChanges()` — executes those changes against the database.

Atlas handles SQLite's structural limitations automatically. When a column type
must change (which SQLite cannot do with `ALTER TABLE`), Atlas generates the
standard workaround: create a new table, copy data, drop the old table, rename
the new one. For purely additive changes (new columns, new tables) it emits
simple `ALTER TABLE … ADD COLUMN` statements.

### Declarative workflow for SFGA

```
desired schema  ──(schema.sql from current sflib tag)──► Atlas Schema object
                                                              │
existing .sqlite ──(Atlas InspectSchema)──────────────► Atlas Schema object
                                                              │
                                    SchemaDiff ───────────────┘
                                         │
                                    ApplyChanges  (on a copy of the .sqlite)
```

The "desired" schema comes from the same `schema.sql` that `Create()` already
fetches — no new source of truth is introduced.

### What changes in `sflib`

#### 1. New dependency

```
ariga.io/atlas v1.2.0
```

Already added to `go.mod`. The SQLite driver lives at
`ariga.io/atlas/sql/sqlite`. Works with `modernc.org/sqlite` (pure Go,
no CGO) — verified by integration test.

#### 2. Interface changes (`pkg/sfga/interface.go`)

**Replace `Updater` with `Migrator`:**

```go
type Migrator interface {
    // Migrate brings the archive's schema up to the current version by
    // computing an Atlas diff between the live schema and the desired
    // schema (schema.sql at the current sflib tag) and applying only the
    // necessary changes. It operates on a copy of the database written
    // to dstDir, so the source file is never modified. Returns the
    // Archive pointing at the migrated copy.
    Migrate(dstDir string) (Archive, error)
}
```

The old `Update(emptySfga Archive, withParents bool) error` method is
**removed entirely**. Its two responsibilities are split:
- Schema migration → `Migrate()` (Atlas-based)
- Parent hierarchy building → `AddParents()` (moved to `Enricher`)

**Simplify `Enricher` — move config to `config.Config`:**

Remove `BasionymInferenceConfig` struct from `interface.go`. Move its fields
into `config.Config`:

```go
// config/config.go additions:
type Config struct {
    // ... existing fields ...

    // InferBasionyms enables basionym inference during enrichment.
    InferBasionyms bool

    // SkipBasionymsIfRelationsExist skips inference if BASIONYM relations
    // already exist in the archive.
    SkipBasionymsIfRelationsExist bool

    // CreateOriginalCombinations creates OriginalGenus, OriginalSpecies,
    // etc. relationships in addition to BASIONYM during inference.
    CreateOriginalCombinations bool

    // MigrateOutputDir, if set, saves an additional copy of the migrated
    // SFGA file to this directory. Migration always happens in the cache
    // directory for internal consistency; this option produces a second
    // copy for user convenience (e.g. sharing, archiving). A simple file
    // copy from cache to the output path.
    MigrateOutputDir string
}
```

The `Enricher` interface becomes:

```go
type Enricher interface {
    // InferBasionyms detects and creates basionym relationships by matching
    // stemmed epithets and original authorship across all names in the
    // archive. Behavior controlled by Config.InferBasionyms,
    // Config.SkipBasionymsIfRelationsExist, and
    // Config.CreateOriginalCombinations.
    InferBasionyms(ctx context.Context) error

    // AddParents builds a parent/child hierarchy from flat classification
    // fields (kingdom, phylum, class, etc.) in the taxon table. Works
    // in-place: deletes and re-inserts name, taxon, synonym, and
    // name_relation tables within a transaction. No-op if parent IDs
    // already exist or flat hierarchy is empty.
    AddParents(ctx context.Context) error
}
```

Both methods read their behavior flags from `a.cfg` — consistent with how
`BatchSize`, `NomCode`, `WithParents` already work.

**Updated `Archive` composition:**

```go
type Archive interface {
    arch.Packager
    AccessorSFGA
    Reader
    Writer
    Migrator
    Enricher
}
```

#### 3. New file `internal/isfga/migrate.go`

Core logic:

```go
func (a *isfga) Migrate(dstDir string) (Archive, error) {
    // 1. Copy the existing .sqlite file to dstDir.
    dstPath := filepath.Join(dstDir, filepath.Base(a.dbPath))
    if err := copyFile(a.dbPath, dstPath); err != nil {
        return nil, err
    }

    // 2. Open the copied file with the Atlas SQLite driver.
    db, err := sql.Open("sqlite", dstPath)
    ...
    drv, err := atlassqlite.Open(db)
    ...

    // 3. Build the desired schema by applying schema.sql to a temp
    //    in-memory SQLite database and inspecting it.
    desired, err := a.desiredSchema(ctx, drv)
    ...

    // 4. Inspect the current (copied) database.
    current, err := drv.InspectSchema(ctx, "main", nil)
    ...

    // 5. Compute diff and apply.
    changes, err := drv.SchemaDiff(current, desired)
    ...
    if len(changes) == 0 {
        slog.Info("SFGA schema is already current, no migration needed")
    } else {
        if err := drv.ApplyChanges(ctx, changes); err != nil {
            return nil, err
        }
        slog.Info("SFGA schema migrated", "changes", len(changes))
    }

    // 6. Update the version table to reflect the new schema version.
    //    The version table stores the schema version as `id` (text).
    //    Atlas only changes structure, not data — so we must update this
    //    ourselves after a successful migration.
    _, err = db.Exec(
        "UPDATE version SET id = ?", config.SchemaVersion,
    )
    ...

    // 7. Return a new Archive pointing at the migrated copy.
    result := New(a.cfg)
    result.SetDb(dstPath)
    _, err = result.Connect()
    return result, err
}
```

#### 4. Refactor `addParents` → in-place `AddParents`

`add_parents.go` is refactored so `AddParents(ctx)` works on the current
archive in-place instead of writing to a separate empty archive:

1. Read all NameUsages → build node tree (same as now).
2. Within a SQLite transaction:
   - `DELETE FROM name; DELETE FROM taxon; DELETE FROM synonym;
     DELETE FROM name_relation;`
   - Re-insert enriched NameUsages with parent IDs via `InsertNameUsages`.
3. On error, transaction rolls back — database unchanged.

The `shouldAddParents()` check remains: no-op if parent IDs already exist
or flat hierarchy is empty.

#### 5. Auto-migration in `Fetch()` (`internal/isfga/fetch.go`)

After `setDb()` succeeds, compare versions. If the file's schema is older
than the current schema, automatically call `Migrate()`:

```go
func (a *isfga) Fetch(src, dstDir string) error {
    ...
    err = a.setDb(dstDir)
    ...
    if gnlib.CmpVersion(a.Version(), config.SchemaVersion) == -1 {
        slog.Warn("SFGA schema outdated, migrating",
            "file_version", a.Version(),
            "current_version", config.SchemaVersion,
        )
        migDir := dstDir + "_migrated"
        migrated, err := a.Migrate(migDir)
        if err != nil {
            return fmt.Errorf("auto-migration failed: %w", err)
        }
        _ = a.Close()
        a.dbPath = migrated.DbPath()
        a.db = migrated.Db()

        // If the user requested a persistent copy, save one.
        if a.cfg.MigrateOutputDir != "" {
            outPath := filepath.Join(
                a.cfg.MigrateOutputDir,
                filepath.Base(a.dbPath),
            )
            if err := copyFile(a.dbPath, outPath); err != nil {
                return fmt.Errorf("saving migrated copy: %w", err)
            }
            slog.Info("Migrated SFGA saved", "path", outPath)
        }
    }
    return nil
}
```

Migration always happens in the cache directory — consistent and predictable.
If `MigrateOutputDir` is set, an additional file copy is made to the user's
chosen path. This means **every consumer** of an SFGA file gets a silently
up-to-date schema with no code changes required upstream in `sf`.

#### 6. Expose `SchemaVersion` constant in `config`

Add a `SchemaVersion` string constant to `config/config.go` derived from
`repoTag` so that `Fetch()` and tests can compare against it without
re-deriving it.

#### 7. Delete `update.go`

All of the following become dead code and are removed:
- `Update()`, `getTables()`, `getFields()`, `transferData()`,
  `transferTableData()`, `attachDatabase()`, `detachDatabase()`,
  `transferDataWithParents()`
- `shouldAddParents()` moves to `add_parents.go`
- `update_test.go` is rewritten to test `Migrate()` and `AddParents()`
  separately

### What changes in `sf` (consumer)

- `sf update` command: simplified. No longer needs two archive arguments.
  Calls `Migrate()` on the source, optionally `AddParents(ctx)` on the
  result if `--add-parents` is set.
- `sf from` subcommands: no change — they call `Create()` not `Fetch()`.
- `sf diff`, `sf apply`, `sf to`: no change needed — migration now happens
  transparently inside `Fetch()`.
- Global `--save-migrated <dir>` flag (or config/env): sets
  `Config.MigrateOutputDir`. When migration triggers during any command,
  the user gets a persistent copy of the updated `.sqlite` file at
  their chosen path. Useful for sharing or avoiding repeated migrations
  on the same input file.

### Caller flow (before → after)

```
# Before: sf update --add-parents old.sfga new.sfga
#   → Update(emptyArchive, true)  // entangled migration + parents

# After:
#   1. Fetch() auto-migrates if needed (transparent to all consumers)
#   2. AddParents(ctx) enriches in-place (explicit, if requested)
```

## Compatibility Semantics

`IsCompatible(version)` currently returns `true` when the file's version is
`>=` the provided version (i.e. file is newer or equal). This is backwards: we
want to migrate when the file's version is **older** than the current schema.

The call site in `Fetch()` should be:

```go
// migrate if file version is older than current schema
if gnlib.CmpVersion(a.Version(), config.SchemaVersion) == -1 {
    // run migration
}
```

This avoids changing the `IsCompatible` semantics for existing callers.

## File Layout After Change

```
internal/isfga/
    isfga.go          — no change
    create.go         — no change
    fetch.go          — add auto-migration call after setDb()
    schema.go         — no change
    update.go         — DELETED
    update_test.go    — REWRITTEN → migrate_test.go + add_parents_test.go
    add_parents.go    — refactored: in-place AddParents(ctx), absorbs
                        shouldAddParents() from update.go
    migrate.go        — NEW: Atlas-based Migrate() implementation
    migrate_test.go   — EXISTS: Atlas + modernc pure-Go verification test
    infer_basionyms.go — updated: reads config from a.cfg instead of
                         BasionymInferenceConfig parameter
pkg/sfga/
    interface.go      — Updater → Migrator, BasionymInferenceConfig removed,
                        Enricher gains AddParents, InferBasionyms simplified
config/
    config.go         — add SchemaVersion constant, add basionym config
                        fields, add Option funcs for new fields
```

## Resolved Questions

1. **SQLite driver compatibility**: VERIFIED — Atlas `sql/sqlite` package is
   driver-agnostic, working purely through `database/sql` interfaces. It does
   not import `mattn/go-sqlite3`. Passing a `*sql.DB` opened with
   `modernc.org/sqlite` via `atlassqlite.Open(db)` works correctly with
   `CGO_ENABLED=0`. Integration test in `internal/isfga/migrate_test.go`
   confirms InspectSchema, SchemaDiff, and ApplyChanges all work. The project
   stays pure Go.

2. **`--add-parents` decoupling**: RESOLVED — `AddParents` moves to `Enricher`
   interface, works in-place within a transaction. No longer coupled to schema
   migration. `update.go` is deleted entirely.

3. **Enricher configuration**: RESOLVED — `BasionymInferenceConfig` struct
   removed. Its fields move to `config.Config` alongside existing flags like
   `WithParents`. Both `InferBasionyms` and `AddParents` read behavior from
   `a.cfg`.

## Open Questions

1. **Read-only source files**: `Migrate()` always works on a copy, so the
   user's original file is never modified. The migrated copy lives in the
   cache directory (`CacheDir/…`). This is consistent with the existing
   pattern where `Fetch()` already extracts into a cache subdirectory.

2. **Migration failure on breaking changes**: if a future schema version
   removes a column that has a `NOT NULL` constraint and no default, Atlas
   will produce a table-rebuild migration. This is safe for SQLite but may
   lose data in the dropped column. Document this in the `Migrate()` godoc
   and log a warning listing the dropped columns.
