package sflib

type SFGA interface {
	Packager
}

type CoLDP interface {
	Packager
}

type DwCA interface {
	Packager
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
