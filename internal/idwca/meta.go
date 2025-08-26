package idwca

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"

	"github.com/sfborg/sflib/pkg/arch"
	dwca "github.com/sfborg/sflib/pkg/dwca"
)

// getMeta reads metadata file and returns dwca.Meta struct. In case if
// something went wrong, it returns an error.
func getMeta(metaPath string) (*dwca.Meta, error) {
	r, err := os.Open(metaPath)
	if err != nil {
		return nil, &arch.ErrFileOpen{Path: metaPath, Err: err}
	}
	defer r.Close()

	var res dwca.Meta
	// Decode XML directly from the file reader to avoid reading the entire file into memory first.
	decoder := xml.NewDecoder(r)
	err = decoder.Decode(&res)
	if err != nil {
		return nil, &arch.ErrMetaDecoder{Err: err}
	}

	// Helper function to parse string indices into integers.
	// It handles empty strings by returning -1 and wraps strconv.Atoi errors with more context.
	parseIntIdx := func(indexStr string, fieldName string) (int, error) {
		if indexStr == "" {
			return -1, nil
		}
		val, err := strconv.Atoi(indexStr)
		if err != nil {
			return 0, fmt.Errorf("cannot convert %s index '%s': %w", fieldName, indexStr, err)
		}
		return val, nil
	}

	// Parse Core ID index
	res.Core.ID.Idx, err = parseIntIdx(res.Core.ID.Index, "meta core ID")
	if err != nil {
		return nil, err
	}

	// Parse Core Fields indices
	for i := range res.Core.Fields {
		res.Core.Fields[i].Idx, err = parseIntIdx(res.Core.Fields[i].Index, "meta core field")
		if err != nil {
			return nil, err
		}
	}

	// Parse Extension CoreID and Fields indices
	for i := range res.Extensions {
		res.Extensions[i].CoreID.Idx, err = parseIntIdx(res.Extensions[i].CoreID.Index, "meta ext ID")
		if err != nil {
			return nil, err
		}
		for j := range res.Extensions[i].Fields {
			res.Extensions[i].Fields[j].Idx, err = parseIntIdx(res.Extensions[i].Fields[j].Index, "meta ext field")
			if err != nil {
				return nil, err
			}
		}
	}

	return &res, nil
}
