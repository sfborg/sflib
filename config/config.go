package config

import (
	"github.com/gnames/gnfmt"
	"github.com/gnames/gnlib/ent/nomcode"
)

var (
	// repoURL is the URL to the SFGA schema repository.
	repoURL = "https://github.com/sfborg/sfga"

	// repoTag of the sfga repo to get correct schema version.
	repoTag = "v0.5.1"

	// schemaHash is the sha256 sum of the correponding schema version.
	schemaHash = "dd7a806e1384d"

	// SchemaVersion is the desired SFGA schema version.
	SchemaVersion = repoTag

	// RepoMinVersion is the oldest SFGA schema version that sflib can migrate.
	// Archives below this version must be brought up to RepoMinVersion using an
	// older sflib release before further migration is possible.
	RepoMinVersion = "v0.5.1"
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

	// LocalSchemaPath is an optional path to a local schema.sql file.
	// If set, it will be used instead of fetching from GitRepo.
	LocalSchemaPath string

	// NomCode represents nomenclatural code relevant for the dataset.
	NomCode nomcode.Code

	// BadRow sets how to process rows with wrong number of fields in CSV
	// files.
	BadRow gnfmt.BadRow

	// WithQuotes sets CSV reader to use `"` as quote. When it is true,
	// RFC-based CSV reader is used, even if delimiter is tab or pipe.
	WithQuotes bool

	// WithParents flag can be set to true when a dataset with a flat
	// hierachy needs conversion to a parent/child hierachy.
	// All IDs generated during unflattening will have 'sf-' prefix.
	WithParents bool

	// BatchSize tells how many elements (rows, structs) to deal with in
	// a chunk of data.
	BatchSize int

	// JobsNum conveys the number of concurrent jobs to run, where it is
	// needed.
	JobsNum int

	// InferBasionyms enables basionym inference during enrichment.
	InferBasionyms bool

	// SkipBasionymsIfRelationsExist skips basionym inference if BASIONYM
	// relations already exist in the archive.
	SkipBasionymsIfRelationsExist bool

	// CreateOriginalCombinations creates OriginalGenus, OriginalSpecies, etc.
	// relationships in addition to BASIONYM during inference.
	CreateOriginalCombinations bool

	// MigrateOutputDir, if set, saves a copy of the migrated SFGA file to
	// this directory after auto-migration in Fetch().
	MigrateOutputDir string
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

func OptNomCode(code nomcode.Code) Option {
	return func(c *Config) {
		c.NomCode = code
	}
}

func OptLocalSchemaPath(path string) Option {
	return func(c *Config) {
		c.LocalSchemaPath = path
	}
}

func OptWithParents(b bool) Option {
	return func(c *Config) {
		c.WithParents = b
	}
}

func OptInferBasionyms(b bool) Option {
	return func(c *Config) {
		c.InferBasionyms = b
	}
}

func OptSkipBasionymsIfRelationsExist(b bool) Option {
	return func(c *Config) {
		c.SkipBasionymsIfRelationsExist = b
	}
}

func OptCreateOriginalCombinations(b bool) Option {
	return func(c *Config) {
		c.CreateOriginalCombinations = b
	}
}

func OptMigrateOutputDir(dir string) Option {
	return func(c *Config) {
		c.MigrateOutputDir = dir
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
		NomCode:   nomcode.Unknown,
		BatchSize: 50_000,
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res
}
