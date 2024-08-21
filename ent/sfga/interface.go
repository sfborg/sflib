// package sfga profides methods for getting the right version of the SFGA
// schema from its github.com/sfborg/sfga repository.
package sfga

// SFGA provides methods to operate SFGA schema.
type SFGA interface {
	// FetchSchema returns SFGA schema according to provided GitRepo.
	// If something went wrong, or the sha256 does not match downloaded
	// schema, it returns an error.
	FetchSchema() ([]byte, error)

	// Clean removes the temporary path for the schema. Returns an error
	// if something went wrong.
	Clean() error

	// GitRepo() returns GitRepo data of the SFGA instance.
	GitRepo() GitRepo

	// Path returns temporary directory where SFGA schema is downloaded from
	// GitRepo.
	Path() string
}
