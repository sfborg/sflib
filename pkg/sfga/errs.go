package sfga

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
