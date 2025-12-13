package sfga

import (
	"context"
	"database/sql"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
)

// Schema defines methods for managing the SFGA database schema.
// Specific data required for methods is taken from the configuraion of
// the Schema instance.
type Schema interface {
	// Fetch retrieves the SFGA schema based on the configured Git repository.
	// Returns the schema in bytes, and an error if retrieval fails or the
	// downloaded schema's SHA256 hash doesn't match the expected value.
	Fetch() ([]byte, error)
}

type Archive interface {
	arch.Packager
	AccessorSFGA
	Reader
	Writer
	Updater
	Enricher
}

// BasionymInferenceConfig controls the basionym inference behavior.
type BasionymInferenceConfig struct {
	// SkipIfRelationsExist skips inference if BASIONYM relations already exist
	SkipIfRelationsExist bool

	// CreateOriginalCombinations creates OriginalGenus, OriginalSpecies, etc.
	// relationships in addition to BASIONYM
	CreateOriginalCombinations bool
}

// Enricher provides methods for enriching SFGA data through inference.
type Enricher interface {
	// InferBasionyms detects and creates basionym relationships by matching
	// stemmed epithets and original authorship across all names in the archive.
	// This is useful for archives that don't have explicit basionym relationships.
	// The config controls whether to skip if relations exist and whether to
	// create OriginalCombination relationships.
	InferBasionyms(ctx context.Context, cfg BasionymInferenceConfig) error
}

// Updater provides methods for upgrading SFGA data from old schemas to
// the last one. It does require schemas to be compatible, meaning that the
// new schema only adds information, and does not remove/modify fields.
type Updater interface {
	// Update takes the URL or local path of an old SFGA file and the output
	// path, which will be used for the SFGA file with the lates schema.
	// Optional flag would create a Zip version of the file. It returns error
	// in case if a problem arises.
	Update(oldSfga Archive) error
}

// Reader provides methods to feed data from SFGA tables to a channel with
// the corresponding CoLDP objects.
type Reader interface {
	LoadMeta() (*coldp.Meta, error)
	LoadAuthors(context.Context, chan<- coldp.Author) error
	LoadDistributions(context.Context, chan<- coldp.Distribution) error
	LoadMedia(context.Context, chan<- coldp.Media) error
	LoadNames(context.Context, chan<- coldp.Name) error
	LoadNameUsages(context.Context, chan<- coldp.NameUsage) error
	LoadNameRelationships(context.Context, chan<- coldp.NameRelation) error
	LoadReferences(context.Context, chan<- coldp.Reference) error
	LoadSpeciesEstimates(context.Context, chan<- coldp.SpeciesEstimate) error
	LoadSpeciesInteractions(context.Context, chan<- coldp.SpeciesInteraction) error
	LoadSynonyms(context.Context, chan<- coldp.Synonym) error
	LoadTaxonConceptRelations(context.Context, chan<- coldp.TaxonConceptRelation) error
	LoadTaxonProperties(context.Context, chan<- coldp.TaxonProperty) error
	LoadTaxa(context.Context, chan<- coldp.Taxon) error
	LoadTreatments(context.Context, chan<- coldp.Treatment) error
	LoadTypeMaterials(context.Context, chan<- coldp.TypeMaterial) error
	LoadVernaculars(context.Context, chan<- coldp.Vernacular) error
}

// AccessorSFGA defines methods for establishing and managing a connection to the
// SQLite database associated with the SFGA archive.
type AccessorSFGA interface {
	// Connect establishes a connection to the SQLite database and returns the
	// database handle or an error if the connection fails.
	Connect() (*sql.DB, error)

	// SetDb allows to change database of the archive.
	SetDb(path string)

	// Db returns connector to the database.
	Db() *sql.DB

	// Ping checks if database exists
	Ping() bool

	// Close terminates the database connection.
	Close() error

	// DbPath returns the path to the SFGA database file. If the file is not
	// yet available, it returns an empty string.
	DbPath() string

	// Version returns the version number of the SFGA schema.
	Version() string

	// IsCompatible returns back true if provided version is equal or larger
	// than the version of the database.
	IsCompatible(version string) bool
}

type Writer interface {
	// InsertMeta saves data from *coldp.Meta object to SFGA DB.
	InsertMeta(meta *coldp.Meta) error

	// InsertAuthors saves data from coldp.Author objects to SFGA DB.
	InsertAuthors(data []coldp.Author) error

	// InsertDistributions saves data from coldp.Distribution objects to
	// SFGA DB.
	InsertDistributions(data []coldp.Distribution) error

	// InsertMedia saves data from coldp.Media objects to SFGA DB.
	InsertMedia(data []coldp.Media) error

	// InsertNames saves data from coldp.Name objects to SFGA DB.
	InsertNames(data []coldp.Name) error

	// InsertNameRelations saves data from coldp.NameRelation objects to
	// SFGA Db.
	InsertNameRelations(data []coldp.NameRelation) error

	// InsertNameUsages saves data from coldp.NameUsage objects to
	// SFGA Db.
	InsertNameUsages(data []coldp.NameUsage) error

	// InsertReferences saves data from coldp.Reference objects to
	// SFGA Db.
	InsertReferences(data []coldp.Reference) error

	// InsertSpeciesEstimates saves data from  coldp.SpeciesEstimate objects to
	// SFGA Db.
	InsertSpeciesEstimates(data []coldp.SpeciesEstimate) error

	// InsertSpeciesInteractions saves data from coldp.SpeciesInteraction objects
	// to SFGA Db.
	InsertSpeciesInteractions(data []coldp.SpeciesInteraction) error

	// InsertSynonyms saves data from coldp.Synonym objects
	// to SFGA Db.
	InsertSynonyms(data []coldp.Synonym) error

	// InsertTaxa saves data from coldp.Taxon objects
	// to SFGA Db.
	InsertTaxa(data []coldp.Taxon) error

	// InsertTaxonConceptRelations saves data from coldp.TaxonConcepRelation
	// objects to SFGA Db.
	InsertTaxonConceptRelations(data []coldp.TaxonConceptRelation) error

	// InsertTaxonProperties saves data from coldp.TaxonProperty objects to SFGA
	// Db.
	InsertTaxonProperties(data []coldp.TaxonProperty) error

	// InsertTreatments saves data from coldp.Treatment objects to SFGA Db.
	InsertTreatments(data []coldp.Treatment) error

	// InsertTypeMaterials saves data from coldp.TypeMaterial objects to SFGA Db.
	InsertTypeMaterials(data []coldp.TypeMaterial) error

	// InsertVernaculars saves data from coldp.Vernacular objects to SFGA Db.
	InsertVernaculars(data []coldp.Vernacular) error
}
