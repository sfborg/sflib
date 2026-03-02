# `sflib`

`sflib` is a Go library containing shared functionality for Species File Group
(SFG) projects. It primarily focuses on handling functionality for Species File
Group Archives (SFGAs).

## Overview

`sflib` provides a unified interface for working with biodiversity data across
five archive formats. The [Catalogue of Life Data Package (CoLDP)][coldp] is
the shared internal data model used throughout.

### Supported Archive Formats

| Package      | Format                          | Description                                              |
|--------------|---------------------------------|----------------------------------------------------------|
| `pkg/text`   | Plain text                      | One scientific name per line (UTF-8)                     |
| `pkg/xsv`    | CSV / TSV / PSV                 | Delimited values with automatic DarwinCore header mapping |
| `pkg/coldp`  | CoLDP                           | Catalogue of Life Data Package (full standard)           |
| `pkg/dwca`   | Darwin Core Archive             | DwCA format as used by GBIF                              |
| `pkg/sfga`   | Species File Group Archive      | SQLite-based archive with full read/write support        |

Each format exposes a consistent `Archive` interface with `Load`, `Write`, and
`Export` methods, created via factory functions in the root package:

```go
sflib.NewText(opts...)
sflib.NewXsv(opts...)
sflib.NewColdp(opts...)
sflib.NewDwca(opts...)
sflib.NewSfga(opts...)
```

### CoLDP Data Model (`pkg/coldp`)

The CoLDP package defines the core types shared across all formats:

* **NameUsage** — consolidated record combining Name, Taxon, and Synonym data
* **Name** — scientific name with parsed components and nomenclatural metadata
* **Taxon** — full taxonomic classification (Kingdom through Realm)
* **Reference** — bibliographic references
* **Author** — person/organisation metadata
* **Distribution** — geographic distribution records
* **Vernacular** — common/vernacular names
* **Synonym**, **NameRelation**, **TaxonConceptRelation** — relationship types
* **Treatment**, **TypeMaterial**, **SpeciesEstimate**, **SpeciesInteraction**,
  **TaxonProperty**, **Media** — extended data types

Comprehensive enumerated types are provided for taxonomic rank, status,
nomenclatural status, habitat, sex, and more.

### Species File Group Archive (`pkg/sfga`)

SFGA is an SQLite-based archive format. The package provides:

* **Reader** — load all CoLDP entity types from the database
* **Writer** — insert all CoLDP entity types into the database
* **AccessorSFGA** — connection management, schema version checking, and
  compatibility validation
* **Updater** — migrate an existing SFGA to a newer schema version
* **Enricher** — enrich data, e.g. automatically infer basionym relationships
* **Schema** — fetch the current SFGA schema from the upstream repository

### Darwin Core Archive (`pkg/dwca`)

* **Reader** — read core taxon file plus Vernacular and Distribution extensions,
  with diagnostics and EML metadata access
* **Writer** — write `meta.xml`, `eml.xml`, core taxon TSV, and extension files

### XSV Format (`pkg/xsv`)

Automatically detects field delimiters and maps DarwinCore and CoLDP column
headers to the canonical CoLDP terms, enabling straightforward import of
third-party CSV/TSV exports.

### Name Parsing (`pkg/parser`)

A thin wrapper around [GNparser][gnparser] that provides nomenclatural-code-
aware parsing and concurrent parser pools for high-throughput workloads.

### Configuration

All factory functions accept functional options:

* `OptNomCode` — nomenclatural code (Zoological, Botanical, Bacterial, …)
* `OptJobsNum` — number of concurrent workers (default: 5)
* `OptBatchSize` — records per batch (default: 50,000)
* `OptWithParents` — reconstruct parent/child hierarchy from flat data
* `OptLocalSchemaPath` — use a local SFGA schema file instead of fetching
  from GitHub

[coldp]: https://github.com/CatalogueOfLife/coldp
[gnparser]: https://github.com/gnames/gnparser


## Installation

To use `sflib` in your Go project, run:

```bash
go get github.com/sfborg/sflib
```

## Usage Examples

### Working with Species File Group Archives (SFGA)

This section demonstrates how to interact with SFG Archives using `sflib`.

#### Creating a new SFGA

 ```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.Create(dif)
  ...
  _, err = sfga.Connect()
  ...
  defer sfga.Close()
  fmt.Println("SFGA created successfully.")
}
```

#### Import

This example demonstrates how to import data into an existing SFGA.

```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.SetDb(dbPath) // Set the path to the existing SFGA database
  ...
  err = sfga.Import(filePath) // Import data from filePath
  ...
  _, err = sfga.Connect() // Connect to the database
  ...
  defer sfga.Close() // Close the connection when done
  fmt.Println("Data imported successfully.")
}
```

#### Ad-hoc connection

This example demonstrates how to connect to an existing SFGA database without
creating a new one.

```go
import (
  "fmt"
  "github.com/sfborg/sflib/pkg/sfgaio"
)

func main() {
  sfga := sfgaio.New()
  err := sfga.SetDb(filePath) // Set the path to the existing SFGA database
  ...
  _, err = sfga.Connect() // Connect to the database
  ...
  defer sfga.Close() // Close the connection when done
  fmt.Println("Connected to SFGA successfully.")
}
```

## Testing

As the library modifies file system, running tests in parallel might
create running conditions, that will break some tests. To make sure
running only one thread with tests either use

```sh
make test
```

or run tests with `-p 1` option:

```sh
go test ./... -p 1
```
