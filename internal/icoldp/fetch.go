package icoldp

import (
	"os"

	"github.com/sfborg/sflib/internal/util"
)

func (a *icoldp) Fetch(src, dstDir string) error {
	dlDir, src, err := util.AssureLocal(src)
	if err != nil {
		return err
	}
	if dlDir != "" {
		defer os.RemoveAll(dlDir)
	}

	err = util.ExtractOrCopy(src, dstDir)
	if err != nil {
		return err
	}
	a.rootDir = dstDir

	return nil
}
