package parser

import (
	"slices"
	"sync"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
)

// ParserCode determines the nomenclatural code for parsing based on the
// provided code from configuration and the code associated with a specific row.
// The purpose of the code is to parse subgenus either like authorship of genus
// (Botanical and Cultivar code), or as an infragenus (all other codes).
// The local code of the row is more important than the code in the
// configuration, however unknown row code does not influence  anything                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    jj.
func ParserCode(code, rowCode nomcode.Code) nomcode.Code {
	res := nomcode.Unknown
	botCodes := []nomcode.Code{nomcode.Botanical, nomcode.Cultivars}
	if slices.Contains(botCodes, code) {
		res = nomcode.Botanical
	}
	if slices.Contains(botCodes, rowCode) {
		res = nomcode.Botanical
	} else if rowCode != nomcode.Unknown {
		res = nomcode.Unknown
	}
	return res
}

// Pool creates a pool of gnparser instances for different nomenclatural codes.
// The size of the pool is determined by the number of jobs.
func Pool(jobsNum int) map[nomcode.Code]*sync.Pool {
	res := make(map[nomcode.Code]*sync.Pool)

	opts := []gnparser.Option{
		gnparser.OptWithDetails(true),
	}

	res[nomcode.Unknown] = &sync.Pool{
		// New is used when pool is empty to generate new parser.
		New: func() any {
			return newParser(opts...)
		},
	}

	for range jobsNum {
		res[nomcode.Unknown].Put(newParser(opts...))
	}

	optsBot := append(opts, gnparser.OptCode(nomcode.Cultivars))
	res[nomcode.Botanical] = &sync.Pool{
		New: func() any {
			return newParser(optsBot...)
		},
	}

	for range jobsNum {
		res[nomcode.Botanical].Put(newParser(optsBot...))
	}

	return res
}

// newParser creates a new instance of gnparser with the given options.
func newParser(opts ...gnparser.Option) gnparser.GNparser {
	cfg := gnparser.NewConfig(opts...)
	parser := gnparser.New(cfg)
	return parser
}
