// package sfga profides methods for getting the right version of the SFGA
// schema from its github.com/sfborg/sfga repository.
package sfga

import "database/sql"

// Schema provides methods to operate Schema schema.
type Schema interface {
	// Fetch returns SFGA schema according to provided GitRepo.
	// If something went wrong, or the sha256 does not match downloaded
	// schema, it returns an error.
	Fetch() ([]byte, error)

	// Clean removes the temporary path for the schema. Returns an error
	// if something went wrong.
	Clean() error

	// GitRepo() returns GitRepo data of the SFGA instance.
	GitRepo() GitRepo

	// Path returns temporary directory where SFGA schema is downloaded from
	// GitRepo.
	Path() string
}

// Archive deals with SFGA files, and connects to their database.
type Archive interface {
	// Extract uncopresses SFGA file and places it in cache, ready to be
	// queried.
	Extract() error

	// Create uses cache directory to create SFGA archive.
	Create(path string) error

	// Clean removes cache directory.
	Clean() error
	DB
}

// DB provides connection to SFGA archive SQLite database.
type DB interface {
	FileDB() string
	Connect() (*sql.DB, error)
	Close() error
}
