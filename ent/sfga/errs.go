package sfga

type ErrUnknownExt struct {
	File string
}

func (e *ErrUnknownExt) Error() string {
	return "unknown extension"
}

type ErrFileOpen struct {
	File string
	Err  error
}

func (e *ErrFileOpen) Error() string {
	return e.Err.Error()
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

type ErrTarGzReader struct {
	File string
	Err  error
}

func (e *ErrTarGzReader) Error() string {
	return e.Err.Error()
}

type ErrZipReader struct {
	File string
	Err  error
}

func (e *ErrZipReader) Error() string {
	return e.Err.Error()
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

type ErrEmptyTar struct {
	File string
	Err  error
}

func (e *ErrEmptyTar) Error() string {
	return "empty tar file"
}

type ErrCacheClean struct {
	Dir string
	Err error
}

func (e *ErrCacheClean) Error() string {
	return e.Err.Error()
}

type ErrRepoCacheClean struct {
	Dir string
	Err error
}

func (e *ErrRepoCacheClean) Error() string {
	return e.Err.Error()
}

type ErrRepoClean struct {
	URL string
	Err error
}

func (e *ErrRepoClean) Error() string {
	return e.Err.Error()
}

type ErrRepoTagCheckout struct {
	Tag string
	Err error
}

func (e *ErrRepoTagCheckout) Error() string {
	return e.Err.Error()
}

type ErrDownload struct {
	URL string
	Err error
}

func (e *ErrDownload) Error() string {
	return e.Err.Error()
}

type ErrExtractArchive struct {
	File string
	Err  error
}

func (e *ErrExtractArchive) Error() string {
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

type ErrZipCreate struct {
	File string
	Err  error
}

func (e *ErrZipCreate) Error() string {
	return e.Err.Error()
}
