package icoldp

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gnames/gnsys"
	"github.com/sfborg/sflib/pkg/coldp"
)

func (a *icoldp) DirInfo() error {
	// get all fiels with their paths inside of the archive.
	paths, err := a.getFiles()
	if err != nil {
		return err
	}

	// find directory where data files reside
	dataDir := getDataDir(paths)
	var dt coldp.DataType
	var metaOK bool

	var refPath string
	var refPri int

	for _, v := range paths {
		dir, file, ext := gnsys.SplitPath(v)

		// if the file is metadata file, get information about metadata.
		metaOK = a.checkMeta(v, file, ext, metaOK)

		if dir != dataDir {
			continue
		}

		dt = coldp.NewDataType(file, ext)
		if dt == coldp.UnkownDT {
			continue
		}

		// Reference files can appear in multiple formats; apply precedence
		// csv/tsv > jsonl > json > bib and store only the best under ReferenceDT.
		if dt == coldp.ReferenceDT || dt == coldp.ReferenceJsonDT || dt == coldp.ReferenceBibtexDT {
			if p := refFormatPriority(ext); p > refPri {
				refPath = v
				refPri = p
			}
		} else {
			a.dataPaths[dt] = v
		}
	}

	if refPath != "" {
		a.dataPaths[coldp.ReferenceDT] = refPath
	}

	return nil
}

// refFormatPriority returns a priority score for a reference file extension.
// Higher score wins: csv/tsv/psv=4, jsonl=3, json=2, bib=1.
func refFormatPriority(ext string) int {
	switch strings.ToLower(ext) {
	case ".csv", ".tsv", ".psv":
		return 4
	case ".jsonl":
		return 3
	case ".json":
		return 2
	case ".bib":
		return 1
	}
	return 0
}

// getDataDir returns dataDir where data files are residing.
func getDataDir(paths []string) string {
	dirs := make(map[string]int)
	for _, v := range paths {
		dir, file, ext := gnsys.SplitPath(v)
		if coldp.NewDataType(file, ext) != coldp.UnkownDT {
			dirs[dir]++
		}
	}

	var res string
	var count int
	for k, v := range dirs {
		if v > count {
			count = v
			res = k
		}
	}
	return res
}

func (a *icoldp) checkMeta(path, file, ext string, metaOK bool) bool {
	if metaOK {
		return true
	}

	file = strings.ToLower(file)
	ext = strings.ToLower(ext)

	if !strings.HasPrefix(file, "meta") {
		return false
	}

	a.metaPath = path
	if strings.HasPrefix(ext, ".json") {
		a.metaType = coldp.JSON
	}
	if ext == "yaml" {
		a.metaType = coldp.YAML
	}
	return true
}

func (a *icoldp) getFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(
		a.rootDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return files, nil
}
