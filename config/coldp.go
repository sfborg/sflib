package config

import (
	"github.com/gnames/gnfmt"
) // Config is a configuration object for the CoLDP archive data processing.
type ConfigCoLDP struct {
	// BadRow sets how to process rows with wrong number of fields in CSV
	// files.
	BadRow gnfmt.BadRow

	// WithQuotes sets CSV reader to use `"` as quote. When it is true,
	// RFC-based CSV reader is used, even if delimiter is tab or pipe.
	WithQuotes bool
}

// OptionCoLDP is a function type that allows to standardize how options to
// the configuration are organized.
type OptionCoLDP func(*ConfigCoLDP)

func OptBadRow(br gnfmt.BadRow) OptionCoLDP {
	return func(c *ConfigCoLDP) {
		c.BadRow = br
	}
}

func OptWithQuotes(b bool) OptionCoLDP {
	return func(c *ConfigCoLDP) {
		c.WithQuotes = b
	}
}

func NewCoLDP(opts ...OptionCoLDP) ConfigCoLDP {
	res := ConfigCoLDP{}
	for _, opt := range opts {
		opt(&res)
	}
	return res
}
