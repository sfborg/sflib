package config

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// tag of the sfga repo to get correct schema version.
	repoTag = "v0.3.33"

	// schemaHash is the sha256 sum of the correponding schema version.
	schemaHash = "e0259fd6fb89a9"
)

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

type ConfigSfga struct {
	GitRepo
}

func NewSFGA() ConfigSfga {
	res := ConfigSfga{
		GitRepo: GitRepo{
			URL:          repoURL,
			Tag:          repoTag,
			ShaSchemaSQL: schemaHash,
		},
	}
	return res
}
