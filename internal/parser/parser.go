package parser

import (
	"sync"

	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
)

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

func newParser(opts ...gnparser.Option) gnparser.GNparser {
	cfg := gnparser.NewConfig(opts...)
	parser := gnparser.New(cfg)
	return parser
}
