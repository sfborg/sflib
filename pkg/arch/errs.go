package arch

import (
	"errors"
	"fmt"

	"github.com/gnames/gnsys"
)

// Shared errors //

// ErrDir is returned when the directory is in an unknown state, or is not a
// directory.
type ErrDir struct {
	DirPath string
}

func (e *ErrDir) Error() string {
	return fmt.Sprintf("Directory '%s' is broken", e.DirPath)
}

type ErrDirCreate struct {
	Dir string
	Err error
}

func (e *ErrDirCreate) Error() string {
	return e.Err.Error()
}

type ErrDirRemove struct {
	Dir string
	Err error
}

func (e *ErrDirRemove) Error() string {
	return e.Err.Error()
}

type ErrDirChange struct {
	Src, Dst string
	Err      error
}

func (e *ErrDirChange) Error() string {
	return e.Err.Error()
}

type ErrUnknownExt struct {
	File string
}

func (e *ErrUnknownExt) Error() string {
	return "unknown extension"
}

type ErrFileOpen struct {
	Path string
	Err  error
}

func (e *ErrFileOpen) Error() string {
	return e.Err.Error()
}

// ErrFileNotFound is returned when the file is not found.
type ErrFileNotFound struct {
	File string
}

func (e *ErrFileNotFound) Error() string {
	return fmt.Sprintf("file '%s' not found", e.File)
}

type ErrFileCreate struct {
	File string
	Err  error
}

func (e *ErrFileCreate) Error() string {
	return e.Err.Error()
}

type ErrFileCopy struct {
	Src, Dst string
	Err      error
}

func (e *ErrFileCopy) Error() string {
	return e.Err.Error()
}

type ErrDownload struct {
	URL string
	Err error
}

func (e *ErrDownload) Error() string {
	return e.Err.Error()
}

type ErrImportArchive struct {
	File string
	Err  error
}

func (e *ErrImportArchive) Error() string {
	return e.Err.Error()
}

type ErrZipCreate struct {
	File string
	Err  error
}

func (e *ErrZipCreate) Error() string {
	return e.Err.Error()
}

// End of shared errors //

// SFGA errors //

type ErrRepoClone struct {
	URL string
	Err error
}

func (e *ErrRepoClone) Error() string {
	return e.Err.Error()
}

type ErrSQLiteConnect struct {
	Err error
}

func (e *ErrSQLiteConnect) Error() string {
	return e.Err.Error()
}

type ErrSQLitePragma struct {
	Err error
}

func (e *ErrSQLitePragma) Error() string {
	return e.Err.Error()
}

type ErrSQLiteExec struct {
	Err error
}

func (e *ErrSQLiteExec) Error() string {
	return e.Err.Error()
}

type ErrSQLiteQuery struct {
	Err error
}

func (e *ErrSQLiteQuery) Error() string {
	return e.Err.Error()
}

type ErrSQLiteLoadSQL struct {
	Err error
}

func (e *ErrSQLiteLoadSQL) Error() string {
	return e.Err.Error()
}

type ErrSQLiteCreateBinary struct {
	File string
	Err  error
}

func (e *ErrSQLiteCreateBinary) Error() string {
	return e.Err.Error()
}

type ErrSQLiteCreateSQL struct {
	File string
	Err  error
}

func (e *ErrSQLiteCreateSQL) Error() string {
	return e.Err.Error()
}

// End of SFGA errors //

// DwCA errors //

// ErrCoreRead is returned when reading the core file fails.
type ErrCoreRead struct {
	Err error
}

func (e *ErrCoreRead) Error() string {
	return fmt.Sprintf("reading core file failed: %v", e.Err)
}

type ErrExtensionRead struct {
	Err error
}

// ErrExtensionRead is returned when reading the extension file fails.
func (e *ErrExtensionRead) Error() string {
	return fmt.Sprintf("reading extension file failed: %v", e.Err)
}

// ErrContext is returned when the context is canceled.
type ErrContext struct {
	Err error
}

func (e *ErrContext) Error() string {
	return fmt.Sprintf("context canceled: %v", e.Err)
}

type ErrSaveCSV struct {
	Err error
}

func (e *ErrSaveCSV) Error() string {
	return fmt.Sprintf("saving csv file failed: %v", e.Err)
}

// ErrMetaReader is an error type for reading meta.xml files.
type ErrMetaReader struct {
	// Err is the original error.
	Err error
}

func (e *ErrMetaReader) Error() string {
	return fmt.Sprintf("cannot read: %v", e.Err)
}

// ErrMetaDecoder is an error type for decoding meta.xml files.
type ErrMetaDecoder struct {
	//  Err is the original error.
	Err error
}

func (e *ErrMetaDecoder) Error() string {
	return fmt.Sprintf("cannot decode meta.xml: %v", e.Err)
}

// ErrMetaFileNotFound is returned when the meta.xml file is not found in the
// extract directory.
var ErrMetaFileNotFound = errors.New("meta file not found")

// ErrMultipleMetaFiles is returned when there are multiple meta.xml files in
// the extract directory.
var ErrMultipleMetaFiles = errors.New("multiple meta files found")

// ErrEmlReader is an error type for reading Eml.xml files.
type ErrEmlReader struct {
	// Err is the original error.
	Err error
}

func (e *ErrEmlReader) Error() string {
	return fmt.Sprintf("cannot read eml.xml: %v", e.Err)
}

// ErrEmlDecoder is an error type for decoding Eml.xml files.
type ErrEmlDecoder struct {
	//  Err is the original error.
	Err error
}

func (e *ErrEmlDecoder) Error() string {
	return fmt.Sprintf("cannot decode eml.xml: %v", e.Err)
}

// ErrUnknownArchiveType is returned when the file type is not supported.
type ErrUnknownArchiveType struct {
	gnsys.FileType
}

func (e *ErrUnknownArchiveType) Error() string {
	return fmt.Sprintf("unknown file type: %s", e.FileType)
}

// ErrExtract is returned when the extraction of the DwCA file fails.
type ErrExtract struct {
	Path string
	Err  error
}

func (e *ErrExtract) Error() string {
	return fmt.Sprintf("extracting '%s' failed: %v", e.Path, e.Err)
}

// End of DwCA errors //
