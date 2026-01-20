package parser

import (
	"github.com/gnames/gnlib/ent/nomcode"
	"github.com/gnames/gnparser"
)

// ParserCode determines the nomenclatural code to use for parsing.
// The code of a row has precedence over globally defined code.
func ParserCode(code, rowCode nomcode.Code) nomcode.Code {
	if rowCode == nomcode.Unknown {
		return code
	}
	return rowCode
}

// Pool creates a map with nomenclatural codes as keys, and pools of
// GNparser with the corresponding to the pool settings. The following
// nomenclatural codes are used currently: Unknosn, Zoological,
// Botanical, Cultivars, Bacterial, Virus. Each pool's capacity is
// equal to jobsNum.
func Pool(jobsNum int) map[nomcode.Code]chan gnparser.GNparser {
	res := make(map[nomcode.Code]chan gnparser.GNparser)
	codes := []nomcode.Code{
		nomcode.Unknown,
		nomcode.Zoological,
		nomcode.Botanical,
		nomcode.Cultivars,
		nomcode.Bacterial,
		nomcode.Virus,
	}

	opts := []gnparser.Option{
		gnparser.OptWithDetails(true),
	}

	for _, code := range codes {
		codeOpts := append(opts, gnparser.OptCode(code))
		cfg := gnparser.NewConfig(codeOpts...)
		res[code] = gnparser.NewPool(cfg, jobsNum)
	}
	return res
}
