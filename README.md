# `SFlib`

[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.18879464.svg)](https://doi.org/10.5281/zenodo.18879464)

`SFlib` is a Go library containing shared functionality for Species File Group
(SFG) projects. It primarily focuses on handling functionality for Species File
Group Archives ([SFGA]s).

## Overview

`SFlib` provides a unified interface for working with biodiversity data across
several archive formats. It allows converting any supported format into
any other one using SFGA as an intermediary.

### Supported Archive Formats

| Package      | Format            | Description                             |
|--------------|-------------------|-----------------------------------------|
| `pkg/sfga`   | [SFGA]            | SQLite-based lossless interchange format |
| `pkg/coldp`  | [CoLDP]           | Catalogue of Life Data Package          |
| `pkg/dwca`   | [DwCA]            | Darwin Core Archive                     |
| `pkg/xsv`    | [CSV] / TSV / PSV | Delimited values files with DwC headers |
| `pkg/text`   | Plain text        | One scientific name per line (UTF-8)    |

Each format exposes a consistent `Archive` interface with `Fetch`, `Create`, and
`Export` methods, created via factory functions in the root package:

```go
sflib.NewText(opts...)
sflib.NewXsv(opts...)
sflib.NewColdp(opts...)
sflib.NewDwca(opts...)
sflib.NewSfga(opts...)
```

### Species File Group Archive (`pkg/sfga`)

[SFGA] is an SQLite-based archive format based heavily on CoLDP standard. SFGA
is designed as a lossless interchange format, preserving all information from
any format it converts to or from.

[Packager API]

[SFGA Interface API]


### CoLDP Data Model (`pkg/coldp`)

The [CoLDP] package provides base models for both SFGA and CoLDP.

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

[Packager API]

[CoLDP Interface API]

### Darwin Core Archive (`pkg/dwca`)

[DwCA] package allows to convert to and from Darwin Core Archives.

### XSV Format (`pkg/xsv`)

XSV package reads and writes comma-separated, tab-separated and
pipe-separated files (CSV,TSV,PSV files). It automatically detects field
delimiters and maps DarwinCore and CoLDP column headers to SFGA.

[Packager API]

[XSV Interface API]

### Text Format (`pkg/text`)

Text package reads and writes a plain text file where every line contains
one scientific name.

[Packager API]

[Text Interface API]

### Configuration

* `OptNomCode` — nomenclatural code (Zoological, Botanical, Bacterial, …)
* `OptJobsNum` — number of concurrent workers (default: 5)
* `OptBatchSize` — records per batch (default: 50,000)
* `OptWithParents` — reconstruct parent/child hierarchy from flat data
* `OptLocalSchemaPath` — use a local SFGA schema file instead of fetching
  from GitHub

## Installation

Requires Go 1.25 or later. No CGO or system libraries are required —
SQLite support is provided by a pure-Go driver.

```bash
go get github.com/sfborg/sflib
```

## Usage Examples

All formats use the same pattern: create an archive with a `New*` factory,
call `Fetch` to load the source into a cache directory, then read records
through a channel.

For more complete real-world usage, see [sf][sf] — the primary tool built
on top of `SFlib`.

[sf]: https://github.com/sfborg/sf

### Load names from a text file

```go
import (
  "context"
  "sync"

  "github.com/gnames/gnlib/ent/nomcode"
  "github.com/sfborg/sflib"
  "github.com/sfborg/sflib/pkg/coldp"
)

a := sflib.NewText()
err := a.Fetch("names.txt", "/tmp/text-cache")

ch := make(chan coldp.NameUsage)
var wg sync.WaitGroup
wg.Add(1)
go func() {
  defer wg.Done()
  for nu := range ch {
    // process nu
  }
}()

err = a.Load(context.Background(), ch, 5, nomcode.Zoological)
close(ch)
wg.Wait()
```

### Convert a DwCA archive to SFGA

```go
import (
  "context"
  "sync"

  "github.com/sfborg/sflib"
  "github.com/sfborg/sflib/pkg/coldp"
)

dwca := sflib.NewDwca()
err := dwca.Fetch("archive.zip", "/tmp/dwca-cache")

sfga := sflib.NewSfga()
err = sfga.Create("/tmp/sfga-cache")
_, err = sfga.Connect()
defer sfga.Close()

ch := make(chan coldp.Data)
var wg sync.WaitGroup
wg.Add(1)
go func() {
  defer wg.Done()
  for d := range ch {
    sfga.InsertNameUsages(d.NameUsages)
  }
}()

err = dwca.LoadCore(context.Background(), ch)
close(ch)
wg.Wait()
```

### Convert an SFGA archive to CoLDP

Continuing from the previous example, the SFGA archive can be converted to any
other supported format. This demonstrates the full round-trip: DwCA → SFGA →
CoLDP.

```go
import (
  "context"
  "path/filepath"
  "sync"

  "github.com/sfborg/sflib"
  "github.com/sfborg/sflib/pkg/coldp"
)

sfga := sflib.NewSfga()
err := sfga.Fetch("archive.sqlite.zip", "/tmp/sfga-cache")
_, err = sfga.Connect()
defer sfga.Close()

coldpDir := "/tmp/coldp-cache"
cl := sflib.NewColdp()
err = cl.Create(coldpDir)

ch := make(chan coldp.NameUsage)
ctx := context.Background()
var wg sync.WaitGroup
wg.Add(1)
go func() {
  defer wg.Done()
  coldp.Write(ctx, ch, filepath.Join(coldpDir, "NameUsage.tsv"))
}()

err = sfga.LoadNameUsages(ctx, ch)
close(ch)
wg.Wait()

err = cl.Export("output.zip", true)
```

## Testing

As the library modifies file system, running tests in parallel might
create running conditions, that will break some tests. To make sure
running only one thread with tests either use

```sh
just test
```

or run tests with `-p 1` option:

```sh
go test ./... -p 1
```

## Authors

* [Dmitry Mozzherin]

## Contributors

* [Geoffrey Ower]

## License

Released under [MIT license]

[CSV]: https://www.ietf.org/rfc/rfc4180.txt
[CoLDP Interface API]: pkg/coldp/interface.go
[CoLDP]: https://github.com/CatalogueOfLife/coldp
[Dmitry Mozzherin]: https://github.com/dimus
[DwCA]: https://dwc.tdwg.org/terms
[Geoffrey Ower]: https://github.com/gdower
[MIT license]: LICENSE
[Packager API]: pkg/arch/interface.go
[SFGA Interface API]: pkg/sfga/interface.go
[SFGA]: https://github.com/sfborg/sfga
[Text Interface API]: pkg/text/interface.go
[XSV Interface API]: pkg/xsv/interface.go
