package sfga

type GitRepo struct {
	// URL to the SFGA schema repository.
	URL string

	// Tag is a version tag of the SFGA repository to use.
	Tag string

	// ShaSumSchema is sha256 hash for the content of schema.sql file.
	// If the has is trucated, only trucated part is checked. If it is
	// empty, no check is done.
	ShaSchemaSQL string
}
