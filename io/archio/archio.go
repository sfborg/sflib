package archio

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
)

type archio struct {
	sfgaFilePath string
	cachePath    string
	downloadPath string
	dbPath       string
	sfga.DB
}

func New(sfgaFilePath, cachePath string) (sfga.Archive, error) {
	var err error
	exists, _ := gnsys.FileExists(sfgaFilePath)
	if !exists && !strings.HasPrefix(sfgaFilePath, "http") {
		err = fmt.Errorf("file not found '%s'", sfgaFilePath)
		return nil, err
	}
	res := &archio{
		sfgaFilePath: sfgaFilePath,
		cachePath:    cachePath,
		downloadPath: filepath.Join(cachePath, "download"),
		dbPath:       filepath.Join(cachePath, "db"),
	}
	return res, nil
}

func (a *archio) Extract() error {
	var err error
	err = a.resetCache()
	if err != nil {
		return err
	}

	if strings.HasPrefix(a.sfgaFilePath, "http") {
		dlPath, err := a.download()
		if err != nil {
			return err
		}
		a.sfgaFilePath = dlPath
	}

	err = a.extract()
	if err != nil {
		return err
	}

	return nil
}

func (a *archio) Create(path string) error {
	return nil
}

func (a *archio) Clean() error {
	switch gnsys.GetDirState(a.cachePath) {
	case gnsys.DirAbsent:
		return gnsys.MakeDir(a.cachePath)
	case gnsys.DirEmpty:
		return nil
	case gnsys.DirNotEmpty:
		return gnsys.CleanDir(a.cachePath)
	default:
		return fmt.Errorf("cannot reset %s", a.cachePath)
	}
}
