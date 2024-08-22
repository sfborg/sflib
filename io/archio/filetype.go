package archio

type fileType int

const (
	unknownFile fileType = iota
	zipFile
	tarGzFile
	sqlFile
	sqliteFile
)
