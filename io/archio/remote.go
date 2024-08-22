package archio

import "github.com/gnames/gnsys"

func (a *archio) download() (string, error) {
	path, err := gnsys.Download(a.sfgaFilePath, a.downloadPath, true)
	return path, err
}
