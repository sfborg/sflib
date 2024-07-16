// package sfga profides methods for getting the right version of the SFGA
// schema from its github.com/sfborg/sfga repository.
package sfga

type SFGA interface {
	FetchSchema() ([]byte, error)
}
