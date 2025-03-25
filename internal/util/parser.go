package util

import (
	"sync"

	"github.com/gnames/gnparser"
	"github.com/gnames/gnparser/ent/nomcode"
	"github.com/sfborg/sflib/pkg/coldp"
)

func newParser(opts ...gnparser.Option) gnparser.GNparser {
	cfg := gnparser.NewConfig(opts...)
	return gnparser.New(cfg)
}

func ParserPool(jobsNum int) map[coldp.NomCode]*sync.Pool {
	res := make(map[coldp.NomCode]*sync.Pool)

	opts := []gnparser.Option{
		gnparser.OptWithDetails(true),
	}

	res[coldp.UnknownNomCode] = &sync.Pool{
		New: func() any {
			return newParser(opts...)
		},
	}
	for range jobsNum {
		res[coldp.UnknownNomCode].Put(newParser(opts...))
	}

	optsBot := append(opts, gnparser.OptCode(nomcode.Cultivar))
	res[coldp.Botanical] = &sync.Pool{
		New: func() any {
			return newParser(optsBot...)
		},
	}
	for range jobsNum {
		res[coldp.Botanical].Put(newParser(optsBot...))
	}
	return res
}
