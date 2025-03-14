package arch

// Shared errors //

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
