package idwca

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/sfborg/sflib/pkg/arch"
	dwca "github.com/sfborg/sflib/pkg/dwca"
)

func getMeta(metaPath string) (*dwca.Meta, error) {
	r, err := os.Open(metaPath)
	if err != nil {
		return nil, &arch.ErrFileOpen{Path: metaPath, Err: err}
	}
	defer r.Close()

	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, &arch.ErrMetaReader{Err: err}
	}

	var res dwca.Meta
	decoder := xml.NewDecoder(bytes.NewReader(bs))
	err = decoder.Decode(&res)
	if err != nil {
		return nil, &arch.ErrMetaDecoder{Err: err}
	}

	if res.Core.ID.Index == "" {
		res.Core.ID.Idx = -1
	} else {
		res.Core.ID.Idx, err = strconv.Atoi(res.Core.ID.Index)
		if err != nil {
			return nil, fmt.Errorf("cannot convert meta core id: %w", err)
		}
	}

	fs := res.Core.Fields
	for i := range fs {
		if fs[i].Index == "" {
			fs[i].Idx = -1
			continue
		}
		fs[i].Idx, err = strconv.Atoi(fs[i].Index)
		if err != nil {
			return nil, fmt.Errorf("cannot covert meta core field: %w", err)
		}
	}

	for i := range res.Extensions {
		if res.Extensions[i].CoreID.Index == "" {
			res.Extensions[i].CoreID.Idx = -1
			continue
		}
		res.Extensions[i].CoreID.Idx, err = strconv.Atoi(res.Extensions[i].CoreID.Index)
		if err != nil {
			return nil, fmt.Errorf("cannot convert meta ext id: %w", err)
		}
		fs := res.Extensions[i].Fields
		for j := range fs {
			if fs[j].Index == "" {
				fs[j].Idx = -1
				continue
			}
			fs[j].Idx, err = strconv.Atoi(fs[j].Index)
			if err != nil {
				return nil, fmt.Errorf("cannot convert meta ext field: %w", err)
			}
		}
	}

	return &res, nil
}
