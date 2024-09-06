// package sfga provides methods for retrieving and working with the
// appropriate version of the Species File Archive (SFGA) schema from the
// GitHub repository at github.com/sfborg/sfga.
package sfga

import "database/sql"

// Schema defines methods for managing the SFGA database schema.
// Specific data required for methods is taken from the configuraion of
// the Schema instance.
type Schema interface {
	// Fetch retrieves the SFGA schema based on the configured Git repository.
	// Returns the schema in bytes, and an error if retrieval fails or the
	// downloaded schema's SHA256 hash doesn't match the expected value.
	Fetch() ([]byte, error)

	// Clean removes the temporary directory used to store repo with the
	// downloaded schema. Returns an error if the removal process encounters any
	// issues.
	Clean() error

	// GitRepo returns the Git repository information (GitRepo struct)
	// associated with this SFGA instance.
	GitRepo() GitRepo

	// Path returns the temporary directory path where the SFGA schema is
	// downloaded from the Git repository.
	Path() string
}

// Archive provides methods for interacting with SFGA archive files and their
// corresponding database.
type Archive interface {
	// Extract decompresses the SFGA archive file and stores it in a cache
	// directory, making it accessible for querying.
	Extract() error

	// Clean removes the cache directory containing the extracted SFGA archive.
	Clean() error
}

// DB defines methods for establishing and managing a connection to the
// SQLite database associated with the SFGA archive.
type DB interface {
	// Connect establishes a connection to the SQLite database and returns the
	// database handle or an error if the connection fails.
	Connect() (*sql.DB, error)

	// Close terminates the database connection.
	Close() error

	// FileDB returns the path to the SFGA database file. If the file is not
	// yet available, it returns an empty string.
	FileDB() string

	// Export SFGA archive to the outputPath, returns error if export fails.
	// If isBin is true, export binary database, instead of SQL dump. If
	// isZip is true, compress as zip file.
	Export(outputPath string, isBin, isZip bool) error

	// Version returns the version number of the SFGA schema.
	Version() string
}
