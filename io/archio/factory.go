package archio

import (
	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
)

func (a *archio) resetCache() error {
	err := a.Clean()
	if err != nil {
		return &sfga.ErrCacheClean{Dir: a.cachePath, Err: err}
	}

	err = gnsys.MakeDir(a.downloadPath)
	if err != nil {
		return err
	}

	err = gnsys.MakeDir(a.dbPath)
	if err != nil {
		return err
	}

	return nil
}
