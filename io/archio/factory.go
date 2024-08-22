package archio

import (
	"github.com/gnames/gnsys"
)

func (a *archio) resetCache() error {
	err := a.Clean()
	if err != nil {
		return err
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
