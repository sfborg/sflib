package config

import "github.com/sfborg/sflib/pkg/sfga"

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// tag of the sfga repo to get correct schema version.
	repoTag = "v0.3.31"

	// schemaHash is the sha256 sum of the correponding schema version.
	schemaHash = "944e70cb8486fd"
)

type ConfigSFGA struct {
	sfga.GitRepo
}

func NewSFGA() ConfigSFGA {
	res := ConfigSFGA{
		GitRepo: sfga.GitRepo{
			URL:          repoURL,
			Tag:          repoTag,
			ShaSchemaSQL: schemaHash,
		},
	}
	return res
}
