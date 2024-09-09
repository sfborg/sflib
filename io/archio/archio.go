package archio

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/ent/sfga"
)

type archio struct {
	// sfgaFilePath is the path to SFGArchive file. The file can be
	// compressed, SQL dump or binary file of SQLite database.
	sfgaFilePath string

	// cachePath is the root of temporary directory to place working copy of
	// SFGArchive. Other paths (dbPath and downloadPath are children of
	// cachePath.
	cachePath string

	// downloadPath is a temporary directory where remote SFGArchive
	// would be downloaded for further processing.
	downloadPath string

	// dbPath is the path where from which SQLite file of SFGArchive script is
	// accessed.
	dbPath string
}

// New creates an instance that is responsible to deal with
// handling SFGArchive on file system level. It does not interact
// with the archive on the SQLite database level.
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

// Extract decompresses SFGArchive (if needed).
func (a *archio) Extract() error {
	var err error
	err = a.resetCache()
	if err != nil {
		return err
	}

	if strings.HasPrefix(a.sfgaFilePath, "http") {
		dlPath, err := a.download()
		if err != nil {
			return &sfga.ErrDownload{URL: a.sfgaFilePath, Err: err}
		}
		a.sfgaFilePath = dlPath
	}

	err = a.extract()
	if err != nil {
		return &sfga.ErrExtractArchive{File: a.sfgaFilePath, Err: err}
	}

	return nil
}

// Create uses database of SFGArchive that exists in the cache directory
// and makes a copy of the Archive at the path provided by the user.
func (a *archio) Create(path string) error {
	return nil
}

// Clean empties cachePath, or creates it if necessary.
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
