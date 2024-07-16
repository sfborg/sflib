package sfga

type GitHubRepo struct {
	// URL to the SFGA schema repository.
	URL string

	// Path is a temporary location to schema files downloaded from GitHub.
	Path string

	// Tag is a version tag of the SFGA repository to use.
	Tag string
}
