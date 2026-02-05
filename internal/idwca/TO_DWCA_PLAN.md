# Plan: Create DwCA from SFGA

## Goal

Implement SFGA-to-DwCA conversion so that data stored in an SFGA
SQLite database can be exported as a standard Darwin Core Archive
(a zip containing meta.xml, eml.xml, and TSV data files).

## Data Flow

```
SFGA DB
  |
  +-- sfga.Reader.LoadMeta()          -> coldp.Meta
  |                                        |
  |                                        +-> dwca.EML   (eml.xml)
  |                                        +-> dwca.Meta  (meta.xml)
  |
  +-- sfga.Reader.LoadNameUsages()    -> coldp.NameUsage  -> Taxon.tsv
  +-- sfga.Reader.LoadVernaculars()   -> coldp.Vernacular -> VernacularName.tsv
  +-- sfga.Reader.LoadDistributions() -> coldp.Distribution -> Distribution.tsv
  |
  +-- Export() -> zip
```

## Changes Required

### 1. EML Struct Additions (`pkg/dwca/eml.go`)

The current EML struct is incomplete for round-trip support. The
following sections present in real-world EML files (e.g.
`testdata/dwca/eml/medium.xml`) are not captured.

#### 1a. AdditionalMetadata

The `<additionalMetadata>` section holds GBIF-specific fields
(citation, logo, dateStamp). Add to `EML`:

```go
type EML struct {
    // ... existing fields ...
    AdditionalMetadata *AdditionalMetadata `xml:"additionalMetadata"`
}

type AdditionalMetadata struct {
    Metadata GBIFMetadata `xml:"metadata>gbif"`
}

type GBIFMetadata struct {
    DateStamp       string `xml:"dateStamp"`
    HierarchyLevel  string `xml:"hierarchyLevel"`
    Citation        string `xml:"citation"`
    ResourceLogoURL string `xml:"resourceLogoUrl"`
}
```

#### 1b. Dataset Distribution URL

The `<distribution>` element inside `<dataset>` provides the
online URL for the dataset itself (distinct from the Distribution
extension). Add to `Dataset`:

```go
type Dataset struct {
    // ... existing fields ...
    Distribution *DatasetDistribution `xml:"distribution"`
}

type DatasetDistribution struct {
    Scope string `xml:"scope,attr,omitempty"`
    URL   OnlineURL `xml:"online>url"`
}

type OnlineURL struct {
    Function string `xml:"function,attr,omitempty"`
    Value    string `xml:",chardata"`
}
```

#### 1c. TaxonomicCoverage

Add to `Coverage`:

```go
type Coverage struct {
    // ... existing fields ...
    TaxonomicCoverage *TaxonomicCoverage `xml:"taxonomicCoverage"`
}

type TaxonomicCoverage struct {
    GeneralTaxonomicCoverage string                    `xml:"generalTaxonomicCoverage"`
    Classifications          []TaxonomicClassification `xml:"taxonomicClassification"`
}

type TaxonomicClassification struct {
    TaxonRankName  string `xml:"taxonRankName"`
    TaxonRankValue string `xml:"taxonRankValue"`
}
```

#### 1d. Missing fields on existing structs

- `Creator`: add `Address` and `OnlineURL` fields (present in
  medium.xml but not in the struct).
- `Contact`: add `OnlineURL` field.
- `AssociatedParty`: add `ElectronicMailAddress` field.

### 2. DwC Term URI Mapping (`pkg/dwca/terms.go`, new file)

The meta.xml uses Darwin Core term URIs. A mapping from CoLDP
field names to DwC URIs is needed.

```go
// DwC namespace prefixes
const (
    DwcNS  = "http://rs.tdwg.org/dwc/terms/"
    DcNS   = "http://purl.org/dc/terms/"
    GbifNS = "http://rs.gbif.org/terms/1.0/"
)

// Core field mapping: coldp header -> DwC term URI
var CoreTerms = map[string]string{
    "col:id":                  DwcNS + "taxonID",
    "col:parentId":            DwcNS + "parentNameUsageID",
    "col:basionymId":          DwcNS + "originalNameUsageID",
    "col:taxonomicStatus":     DwcNS + "taxonomicStatus",
    "col:scientificName":      DwcNS + "scientificName",
    "col:authorship":          DwcNS + "scientificNameAuthorship",
    "col:rank":                DwcNS + "taxonRank",
    "col:genericName":         DwcNS + "genericName",
    "col:infragenericEpithet": DwcNS + "infragenericEpithet",
    "col:specificEpithet":     DwcNS + "specificEpithet",
    "col:infraspecificEpithet":DwcNS + "infraspecificEpithet",
    "col:cultivarEpithet":     DwcNS + "cultivarEpithet",
    "col:code":                DwcNS + "nomenclaturalCode",
    "col:nameStatus":          DwcNS + "nomenclaturalStatus",
    "col:kingdom":             DwcNS + "kingdom",
    "col:phylum":              DwcNS + "phylum",
    "col:class":               DwcNS + "class",
    "col:order":               DwcNS + "order",
    "col:family":              DwcNS + "family",
    "col:genus":               DwcNS + "genus",
    "col:subgenus":            DwcNS + "subgenus",
    "col:remarks":             DwcNS + "taxonRemarks",
    "col:link":                DcNS  + "references",
    "col:modified":            DcNS  + "modified",
}
```

Fields with the `sf:` prefix (e.g. `sf:speciesId`, `sf:realmId`)
have no standard DwC terms. Two options:

- **Drop them** from the DwCA core (simplest, standard-compliant).
- **Include as custom terms** under a project-specific namespace.

Recommendation: drop `sf:` fields. DwCA consumers will not
understand them. The data can be recovered by going back to SFGA
or CoLDP.

Decision from @dimus: drop them

### 3. coldp.Meta to EML Conversion (`pkg/dwca/eml_convert.go`, new file)

A function to build an EML from coldp.Meta:

```go
func MetaFromCoLDP(m *coldp.Meta) *EML
```

Field mapping:

| coldp.Meta field  | EML location                                      |
|-------------------|---------------------------------------------------|
| Title             | Dataset.Title                                     |
| Description       | Dataset.Abstract.Para                             |
| DOI               | Dataset.AlternativeIdentifier.Value               |
| Issued            | Dataset.PubDate                                   |
| License           | Dataset.IntellectualRights.Para                   |
| URL               | Dataset.Distribution.URL                          |
| Keywords          | Dataset.KeywordSets (single set)                  |
| Contact           | Dataset.Contacts (Actor -> Contact)               |
| Creators          | Dataset.Creators (Actor -> Creator)               |
| Editors           | Dataset.AssociatedParties (role = "editor")        |
| Contributors      | Dataset.AssociatedParties (role = "contributor")   |
| GeographicScope   | Dataset.Coverage.GeographicCoverage.Description   |
| TaxonomicScope    | Dataset.Coverage.TaxonomicCoverage.General...     |
| TemporalScope     | Dataset.Coverage.TemporalCoverage (parse string)  |
| Citation          | AdditionalMetadata.GBIF.Citation                  |
| Logo              | AdditionalMetadata.GBIF.ResourceLogoURL           |

Actor to EML person/org mapping:

```
Actor.Given  + Actor.Family  -> IndividualName
Actor.Organization           -> OrganizationName
Actor.Email                  -> ElectronicMailAddress
Actor.Country/City           -> Address
```

### 4. Meta Generation (`pkg/dwca/meta_build.go`, new file)

A function to build a `Meta` struct from a set of field headers:

```go
func BuildMeta(
    coreHeaders []string,
    extensions map[string][]string,  // name -> headers
) *Meta
```

This uses the CoreTerms mapping (section 2) to produce `<field>`
elements with proper `index` and `term` attributes. The core
rowType is always `http://rs.tdwg.org/dwc/terms/Taxon`.

Extension rowTypes:

| Extension        | RowType URI                                   |
|------------------|-----------------------------------------------|
| VernacularName   | http://rs.gbif.org/terms/1.0/VernacularName   |
| Distribution     | http://rs.gbif.org/terms/1.0/Distribution     |

Headers that have no DwC URI mapping are skipped (not included
in meta.xml fields).

### 5. Archive Interface Changes (`pkg/dwca/interface.go`)

Split the flat interface into `Reader` and `Writer`
sub-interfaces, consistent with the `sfga.Archive` pattern.

```go
type Archive interface {
    arch.Packager
    Reader
    Writer
}

// Reader groups all methods for consuming an existing DwCA.
type Reader interface {
    // Meta returns a pointer to the archive's Meta object.
    Meta() *Meta
    // EML returns a pointer to the archive's provenance EML object.
    EML() *EML
    // Diagnostics returns diagnostics that ran during loading.
    Diagnostics() *diagn.Diagnostics
    // LoadCore reads the core file and converts rows to
    // coldp.NameUsage and coldp.Reference objects.
    LoadCore(ctx context.Context, ch chan<- coldp.Data) error
    // LoadVernacular reads the vernacular names extension.
    LoadVernacular(
        ctx context.Context,
        idx int,
        ch chan<- []coldp.Vernacular,
    ) error
    // LoadDistribution reads the distribution extension.
    LoadDistribution(
        ctx context.Context,
        idx int,
        ch chan<- coldp.Data,
    ) error
    // CoreSlice returns a slice of rows from the core file.
    CoreSlice(offset, limit int) ([][]string, error)
    // CoreStream streams rows of the core file into a channel.
    CoreStream(
        ctx context.Context,
        chCore chan<- []string,
    ) (int, error)
    // ExtensionSlice returns a slice of rows from an extension.
    ExtensionSlice(index, offset, limit int) ([][]string, error)
    // ExtensionStream streams rows of an extension into a channel.
    ExtensionStream(
        ctx context.Context,
        index int,
        ch chan<- []string,
    ) (int, error)
}

// Writer groups all methods for creating a new DwCA.
// Every method writes its output to files inside the archive's
// rootDir (set by Create). Export then only packages those
// cached files into a zip.
type Writer interface {
    // WriteMeta marshals the Meta struct to meta.xml on disk.
    WriteMeta(*Meta) error
    // WriteEML converts coldp.Meta to EML and writes eml.xml on disk.
    // If meta is nil or has no title, a synthetic EML is generated.
    WriteEML(*coldp.Meta) error
    // WriteCore writes core taxon rows from a NameUsage channel
    // to Taxon.tsv on disk.
    WriteCore(ctx context.Context, ch <-chan coldp.NameUsage) error
    // WriteVernaculars writes vernacular name extension rows
    // to VernacularName.tsv on disk.
    WriteVernaculars(
        ctx context.Context,
        ch <-chan coldp.Vernacular,
    ) error
    // WriteDistributions writes distribution extension rows
    // to Distribution.tsv on disk.
    WriteDistributions(
        ctx context.Context,
        ch <-chan coldp.Distribution,
    ) error
}
```

This keeps the single `Archive` type for callers that need both
directions, while allowing code that only writes (e.g. the SFGA
converter) to accept just `dwca.Writer`.

### 6. Implementation (`internal/idwca/`)

The `idwca` struct already satisfies `Reader` (all read methods
exist). The new work is implementing `Writer` methods and
completing the `arch.Packager` stubs (`Create`, `Export`).

#### 6a. Create (`internal/idwca/create.go`, new file)

```go
func (a *idwca) Create(dir string) error
```

- Set `a.rootDir = dir`.
- Create directory if it does not exist.

#### 6b. WriteMeta / WriteEML

```go
func (a *idwca) WriteMeta(m *dwca.Meta) error
```

- Store `a.meta = m`.
- Marshal `m` to XML and write `meta.xml` into `a.rootDir`.

```go
func (a *idwca) WriteEML(e *dwca.EML) error
```

- Store `a.eml = e`.
- Marshal `e` to XML and write `eml.xml` into `a.rootDir`.

#### 6c. WriteCore (`internal/idwca/write_core.go`, new file)

```go
func (a *idwca) WriteCore(
    ctx context.Context,
    ch <-chan coldp.NameUsage,
) error
```

- Open `Taxon.tsv` in `a.rootDir`.
- Write TSV header (DwC field names without URI prefix).
- For each NameUsage from channel, write a TSV row using a
  subset of fields that have DwC term mappings.
- Use tab separator, no enclosure, UTF-8.

The subset of NameUsage fields to include in the core:

```
taxonID, parentNameUsageID, acceptedNameUsageID,
taxonomicStatus, scientificName, scientificNameAuthorship,
taxonRank, nomenclaturalCode, nomenclaturalStatus,
genericName, infragenericEpithet, specificEpithet,
infraspecificEpithet, cultivarEpithet,
kingdom, phylum, class, order, family, genus, subgenus,
taxonRemarks, references, modified
```

Note: `acceptedNameUsageID` — for synonyms, this is the ID of
the accepted taxon. In SFGA, synonyms point to accepted taxa
via the `synonym` table. When using `LoadNameUsages()`, synonyms
should have their `ParentID` repurposed or a separate field used.
This needs investigation during implementation.

#### 6d. WriteVernaculars (`internal/idwca/write_vern.go`, new file)

TSV with fields: `taxonID`, `vernacularName`, `language`,
`countryCode`, `source`.

#### 6e. WriteDistributions (`internal/idwca/write_distr.go`, new file)

TSV with fields: `taxonID`, `occurrenceStatus`, `locationID`,
`locality`, `countryCode`, `source`.

#### 6f. Export (`internal/idwca/export.go`, new file)

```go
func (a *idwca) Export(outputPath string, isZip bool) error
```

By this point all files (meta.xml, eml.xml, Taxon.tsv,
VernacularName.tsv, Distribution.tsv) already exist in
`a.rootDir`, written by the Writer methods. Export only
packages them:

- Walk `a.rootDir` and collect all files.
- If `isZip`: create a zip archive at `outputPath`.
- Otherwise: create a tar.gz archive at `outputPath`.

#### 6g. XML Marshaling Concern

Go's `encoding/xml` has limited namespace support. The EML root
element requires multiple namespace declarations:

```xml
<eml:eml xmlns:eml="eml://ecoinformatics.org/eml-2.1.1"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="..."
    packageId="..." system="http://gbif.org" ...>
```

Options:
- Write a thin XML template for the EML wrapper and inject the
  marshaled `<dataset>` and `<additionalMetadata>` inside it.
- Use `xml.Name{Space: ..., Local: ...}` for the root element
  and add namespace attrs explicitly.

Recommendation: use a template for the EML envelope, marshal the
inner sections (`<dataset>`, `<additionalMetadata>`) with
`encoding/xml`, and stitch them together. This gives full control
over namespace declarations.

For `meta.xml`, marshaling is simpler — only one namespace
(`http://rs.tdwg.org/dwc/text/`).

### 7. Vernacular / Distribution Term Mappings

Similar to core terms, extension fields need URI mappings:

**VernacularName extension:**

```
taxonID        -> http://rs.tdwg.org/dwc/terms/taxonID
vernacularName -> http://rs.tdwg.org/dwc/terms/vernacularName
language       -> http://purl.org/dc/terms/language
countryCode    -> http://rs.tdwg.org/dwc/terms/countryCode
source         -> http://purl.org/dc/terms/source
```

**Distribution extension:**

```
taxonID          -> http://rs.tdwg.org/dwc/terms/taxonID
occurrenceStatus -> http://rs.tdwg.org/dwc/terms/occurrenceStatus
locationID       -> http://rs.tdwg.org/dwc/terms/locationID
locality         -> http://rs.tdwg.org/dwc/terms/locality
countryCode      -> http://rs.tdwg.org/dwc/terms/countryCode
source           -> http://purl.org/dc/terms/source
```

## Implementation Order

1. EML struct additions (section 1) — no breaking changes.
2. DwC term mapping constants (section 2).
3. coldp.Meta -> EML conversion (section 3).
4. Meta builder (section 4).
5. Interface additions (section 5).
6. Write methods + Create + Export (section 6).
7. Tests using existing testdata archives for round-trip
   verification.

## Open Questions

- **sf: fields**: Drop from DwCA, or keep under custom namespace?
- **Synonym handling**: `LoadNameUsages` returns both accepted taxa
  and synonyms. In DwCA, synonyms use `acceptedNameUsageID` in
  the core file (not `parentNameUsageID`). Need to verify how
  `coldp.NameUsage.ParentID` is populated for synonyms coming out
  of SFGA — it may point to the accepted taxon already, or we may
  need to join with the `synonym` table explicitly.
- **Additional extensions**: Should we support SpeciesProfile
  (extinct, marine, freshwater, terrestrial) as seen in the CoL
  DwCA? This can be deferred.
- **Archive format**: Default to zip (GBIF standard) or tar.gz?
  The `isZip` parameter on Export already handles this.
