package idwca

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"

	"github.com/sfborg/sflib/pkg/arch"
	dwca "github.com/sfborg/sflib/pkg/dwca"
)

func getEML(emlPath string) (*dwca.EML, error) {
	r, err := os.Open(emlPath)
	if err != nil {
		return nil, &arch.ErrFileOpen{Path: emlPath, Err: err}
	}

	defer r.Close()
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, &arch.ErrEmlReader{Err: err}
	}

	var res dwca.EML
	decoder := xml.NewDecoder(bytes.NewReader(bs))
	err = decoder.Decode(&res)
	if err != nil {
		return nil, &arch.ErrEmlDecoder{Err: err}
	}

	return &res, nil
}
