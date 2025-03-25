package arch

// Packager defines an interface for managing archive operations, including
// fetching, creating, and exporting archives.
type Packager interface {
	// Fetch retrieves data from a source path and stores it at a destination
	// path.
	Fetch(src, dst string) error

	// Create generates an empty archive at the given directory.
	Create(dir string) error

	// Export writes an archive from the cache to the specified output path. The
	// export can be compressed as a zip file if isZip is true.
	Export(outputPath string, isZip bool) error
}
