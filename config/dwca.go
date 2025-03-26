package config

import "github.com/gnames/gnfmt"

type ConfigDwCA struct {
	BatchSize int
	BadRow    gnfmt.BadRow
}

type OptionDwCA func(*ConfigDwCA)

func OptBatchSize(i int) OptionDwCA {
	return func(cfg *ConfigDwCA) {
		cfg.BatchSize = i
	}
}

func NewDwca(opts ...OptionDwCA) ConfigDwCA {
	res := ConfigDwCA{
		BatchSize: 50_000,
		BadRow:    gnfmt.ProcessBadRow,
	}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
