// package sfga provides methods for retrieving and working with the
// appropriate version of the Species File Archive (SFGA) schema from the
// GitHub repository at github.com/sfborg/sfga.
package sfga

import (
	"database/sql"

	"github.com/gnames/coldp/ent/coldp"
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

// Archive represents SFGA archive
type Archive interface {
	Packager
	Accessor
	CoLDPInserter
}

// Packager provides methods for interacting with SFGA archive packages.
// It can extract files or create a new package from files.
type Packager interface {
	// Extract decompresses the SFGA archive file and stores it in a cache
	// directory, making it accessible for querying.
	Import(src, dst string) error

	// Create a new SFGA file of a specific version.
	Create(dir string) error

	// Export SFGA archive from cache to the outputPath, returns error if export
	// fails. If isBin is true, export binary database, instead of SQL dump. If
	// isZip is true, compress as zip file.
	Export(outputPath string, isZip bool) error
}

// Accessor defines methods for establishing and managing a connection to the
// SQLite database associated with the SFGA archive.
type Accessor interface {
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

type CoLDPInserter interface {
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
