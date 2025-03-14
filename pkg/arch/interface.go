package arch

type DwCA interface {
	Packager
}

// Packager provides methods for interacting with an archive packages.
// It can extract files or create a new package.
type Packager interface {
	// Import decompresses the an archive file and stores it in a cache
	// directory, making it accessible for querying.
	Import(src, dst string) error

	// Create a cached version of an archive from scratch.
	Create(dir string) error

	// Export an archive from cache to the outputPath, returns error if export
	// fails. If isBin is true, export binary database, instead of SQL dump. If
	// isZip is true, compress as zip file.
	Export(outputPath string, isZip bool) error
}
