package sflib

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

type ErrExtractArchive struct {
	File string
	Err  error
}

func (e *ErrExtractArchive) Error() string {
	return e.Err.Error()
}

type ErrZipCreate struct {
	File string
	Err  error
}

func (e *ErrZipCreate) Error() string {
	return e.Err.Error()
}
