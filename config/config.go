package config

import "github.com/sfborg/sflib/ent/sfga"

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// tag of the sfga repo to get correct schema version.
	repoTag = "v0.3.24"

	// schemaHash is the sha256 sum of the correponding schema version.
	schemaHash = "b1db9df2e759f"
)

type Config struct {
	sfga.GitRepo
}

func New() *Config {
	res := Config{
		GitRepo: sfga.GitRepo{
			URL:          repoURL,
			Tag:          repoTag,
			ShaSchemaSQL: schemaHash,
		},
	}
	return &res
}
