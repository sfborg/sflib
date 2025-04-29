package config

import (
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/nomcode"
)

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// repoTag of the sfga repo to get correct schema version.
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

type Config struct {
	// GitRepo is used for initialization of SFGA archive.
	GitRepo

	Code nomcode.Code

	// BadRow sets how to process rows with wrong number of fields in CSV
	// files.
	BadRow gnfmt.BadRow

	// WithQuotes sets CSV reader to use `"` as quote. When it is true,
	// RFC-based CSV reader is used, even if delimiter is tab or pipe.
	WithQuotes bool

	// BatchSize tells how many elements (rows, structs) to deal with in
	// a chunk of data.
	BatchSize int

	// JobsNum conveys the number of concurrent jobs to run, where it is
	// needed.
	JobsNum int
}

type Option func(*Config)

func OptBadRow(br gnfmt.BadRow) Option {
	return func(c *Config) {
		c.BadRow = br
	}
}

func OptWithQuotes(b bool) Option {
	return func(c *Config) {
		c.WithQuotes = b
	}
}

func OptBatchSize(i int) Option {
	return func(cfg *Config) {
		cfg.BatchSize = i
	}
}

func OptJobsNum(i int) Option {
	return func(c *Config) {
		c.JobsNum = i
	}
}

func OptCode(code nomcode.Code) Option {
	return func(c *Config) {
		c.Code = code
	}
}

func New(opts ...Option) Config {
	gitRepo := GitRepo{
		URL:          repoURL,
		Tag:          repoTag,
		ShaSchemaSQL: schemaHash,
	}

	res := Config{
		GitRepo:   gitRepo,
		BadRow:    gnfmt.ProcessBadRow,
		JobsNum:   5,
		Code:      nomcode.Unknown,
		BatchSize: 50_000,
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res
}
