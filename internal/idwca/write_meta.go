package idwca

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sfborg/sflib/pkg/dwca"
)

const metaEnvelopeOpen = xml.Header +
	"<archive xmlns=\"http://rs.tdwg.org/dwc/text/\" metadata=\"eml.xml\">\n"

const metaEnvelopeClose = "</archive>\n"

func (a *idwca) WriteMeta(m *dwca.Meta) error {
	a.meta = m

	bs, err := marshalMeta(m)
	if err != nil {
		return fmt.Errorf("cannot marshal meta.xml: %w", err)
	}

	path := filepath.Join(a.rootDir, "meta.xml")
	if err := os.WriteFile(path, bs, 0644); err != nil {
		return fmt.Errorf("cannot write meta.xml: %w", err)
	}

	return nil
}

func marshalMeta(m *dwca.Meta) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(metaEnvelopeOpen)

	enc := xml.NewEncoder(&buf)
	enc.Indent("  ", "  ")

	if err := enc.EncodeElement(m.Core, xml.StartElement{
		Name: xml.Name{Local: "core"},
	}); err != nil {
		return nil, err
	}

	for _, ext := range m.Extensions {
		if err := enc.EncodeElement(ext, xml.StartElement{
			Name: xml.Name{Local: "extension"},
		}); err != nil {
			return nil, err
		}
	}

	if err := enc.Flush(); err != nil {
		return nil, err
	}

	buf.WriteString("\n")
	buf.WriteString(metaEnvelopeClose)

	return buf.Bytes(), nil
}
