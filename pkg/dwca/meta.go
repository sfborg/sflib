package dwca

import (
	"encoding/xml"
	"path/filepath"
	"strconv"
	"strings"
)

type Meta struct {
	XMLName     xml.Name     `xml:"archive"`
	XMLNameStar xml.Name     `xml:"starArchive"`
	EMLFile     string       `xml:"metadata,attr"`
	Core        *Core        `xml:"core"`
	Extensions  []*Extension `xml:"extension"`
}

// Attr holds the common fields for Core and Extension.
type Attr struct {
	Encoding           string  `xml:"encoding,attr"`
	FieldsTerminatedBy string  `xml:"fieldsTerminatedBy,attr"`
	LinesTerminatedBy  string  `xml:"linesTerminatedBy,attr"`
	FieldsEnclosedBy   string  `xml:"fieldsEnclosedBy,attr"`
	IgnoreHeaderLines  string  `xml:"ignoreHeaderLines,attr"`
	RowType            string  `xml:"rowType,attr"`
	Files              Files   `xml:"files"`
	Fields             []Field `xml:"field"`
}

// Core includes CommonElement and any core-specific fields (like ID).
type Core struct {
	ID ID `xml:"id"`
	*Attr
}

// Extension includes CommonElement and any extension-specific fields (like CoreID).
type Extension struct {
	CoreID CoreID `xml:"coreid"`
	*Attr
}

// Files holds the location of files.
type Files struct {
	// Locations provides path to a file.
	Locations []string `xml:"location"`
}

// ID holds the fields for the Core ID.
type ID struct {
	Index string `xml:"index,attr"`
	Idx   int    `xml:"-"`
	Term  string `xml:"term,attr,omitempty"`
}

// CoreID holds the fields for the CoreID data.
type CoreID struct {
	Index string `xml:"index,attr"`
	Idx   int    `xml:"-"`
}

// Field holds the fields of the data.
type Field struct {
	// Index is the verbatim index of the field.
	Index string `xml:"index,attr"`

	// Idx is the int version of Index.
	Idx int `xml:"-"`

	// Term is the URI of the term.
	Term string `xml:"term,attr"`
}

func (m *Meta) Simplify() *MetaSimple {
	data := &MetaSimple{}
	data.ExtensionsData = make(map[string]ExtensionData)
	data.CoreData = m.Core.toCoreData()

	for _, ext := range m.Extensions {
		file := filepath.Base(ext.Files.Locations[0])
		name := stripExt(file)
		if ext.RowType != "" {
			name = filepath.Base(ext.RowType)
		}
		name = strings.ToLower(name)
		data.ExtensionsData[name] = ext.toExtensionData()
	}
	return data
}

func (c *Core) toCoreData() CoreData {
	var termFull string
	if c.ID.Term != "" {
		termFull = c.ID.Term
	}
	if c.RowType != "" {
		termFull = c.RowType
	}
	var idx int
	idxRes, err := strconv.Atoi(c.ID.Index)
	if err == nil {
		idx = idxRes
	}
	term := filepath.Base(termFull)
	term = strings.ToLower(term)
	coreData := CoreData{
		Index:      idx,
		Locations:  c.Files.Locations,
		TermFull:   termFull,
		Term:       term,
		FieldsData: make(map[string]FieldData),
		FieldsIdx:  make(map[int]FieldData),
	}

	for _, field := range c.Fields {
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)

		if err == nil {
			idx = idxRes
		}

		term := filepath.Base(field.Term)
		term = strings.ToLower(term)
		fd := FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
		coreData.FieldsData[term] = fd
		coreData.FieldsIdx[idx] = fd
	}
	return coreData
}

func (e *Extension) toExtensionData() ExtensionData {
	idx := 0
	idxRes, err := strconv.Atoi(e.CoreID.Index)
	if err == nil {
		idx = idxRes
	}
	extData := ExtensionData{
		CoreIndex:  idx,
		Locations:  e.Files.Locations,
		FieldsData: make(map[string]FieldData),
		FieldsIdx:  make(map[int]FieldData),
	}
	for _, field := range e.Fields {
		term := filepath.Base(field.Term)
		idx = 0
		idxRes, err := strconv.Atoi(field.Index)
		if err == nil {
			idx = idxRes
		}
		fd := FieldData{
			Index:    idx,
			TermFull: field.Term,
			Term:     term,
		}
		extData.FieldsData[term] = fd
		extData.FieldsIdx[idx] = fd
	}
	return extData
}

func stripExt(filename string) string {
	ext := len(filepath.Ext(filename))
	end := len(filename) - ext
	return filename[:end]
}
