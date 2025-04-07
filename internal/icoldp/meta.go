package icoldp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfborg/sflib/pkg/coldp"
	"gopkg.in/yaml.v3"
)

func (a *icoldp) Meta() (*coldp.Meta, error) {
	var err error
	if a.meta != nil {
		return a.meta, nil
	}
	if a.metaPath == "" {
		err = a.DirInfo()
		if err != nil {
			return nil, err
		}
	}

	bs, err := os.ReadFile(a.metaPath)
	if err != nil {
		return nil, err
	}

	var meta coldp.Meta
	if isJSON(a.metaPath) {
		err = json.Unmarshal(bs, &meta)
	} else {
		err = yaml.Unmarshal(bs, &meta)
	}
	if err != nil {
		return nil, err
	}

	a.meta = &meta

	return a.meta, nil
}

func (a *icoldp) WriteMeta(meta *coldp.Meta, path string) error {
	bs, err := json.MarshalIndent(meta, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bs, 0644)
	return nil
}

func isJSON(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	if ext == ".json" {
		return true
	}
	return false
}
