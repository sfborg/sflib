package idwca

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sfborg/sflib/pkg/arch"
	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/sfborg/sflib/pkg/dwca"
)

const emlEnvelopeOpen = xml.Header +
	`<eml:eml xmlns:eml="eml://ecoinformatics.org/eml-2.1.1"
    xmlns:dc="http://purl.org/dc/terms/"
    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="eml://ecoinformatics.org/eml-2.1.1 http://rs.gbif.org/schema/eml-gbif-profile/1.0.1/eml.xsd"
    packageId="%s" system="http://gbif.org" xml:lang="%s">
`

const emlEnvelopeClose = "</eml:eml>\n"

func (a *idwca) WriteEML(m *coldp.Meta) error {
	var e *dwca.EML
	if m == nil || m.Title == "" {
		e = syntheticEML()
	} else {
		e = metaToEML(m)
	}
	a.eml = e

	bs, err := marshalEML(e)
	if err != nil {
		return &arch.ErrEmlDecoder{Err: err}
	}

	path := filepath.Join(a.rootDir, "eml.xml")
	if err := os.WriteFile(path, bs, 0644); err != nil {
		return &arch.ErrFileCreate{File: path, Err: err}
	}

	return nil
}

func syntheticEML() *dwca.EML {
	now := time.Now().Format("2006-01-02")
	return &dwca.EML{
		Lang: "eng",
		Dataset: dwca.Dataset{
			Title:   "Placeholder DwCA (no metadata provided)",
			PubDate: now,
			Abstract: dwca.Abstract{
				Para: "This EML was auto-generated because the source " +
					"archive did not contain dataset metadata.",
			},
			Language: "eng",
		},
	}
}

func metaToEML(m *coldp.Meta) *dwca.EML {
	e := &dwca.EML{
		Lang: "eng",
		Dataset: dwca.Dataset{
			ID:       m.Key,
			Title:    m.Title,
			PubDate:  m.Issued,
			Language: "eng",
			Abstract: dwca.Abstract{
				Para: m.Description,
			},
		},
	}

	// DOI as alternate identifier
	if m.DOI != "" {
		e.Dataset.AlternativeIdentifier = dwca.AltID{
			System: "doi",
			Value:  m.DOI,
		}
	}

	// License
	if m.License != "" {
		e.Dataset.IntellectualRights = &dwca.IntellectualRights{
			Para: m.License,
		}
	}

	// Keywords
	if len(m.Keywords) > 0 {
		kws := make([]dwca.Keyword, len(m.Keywords))
		for i, k := range m.Keywords {
			kws[i] = dwca.Keyword{Value: k}
		}
		e.Dataset.KeywodSets = []dwca.KeywordSet{{Keywords: kws}}
	}

	// Creators
	for _, actor := range m.Creators {
		e.Dataset.Creators = append(e.Dataset.Creators, actorToCreator(actor))
	}

	// Contact
	if m.Contact != nil {
		e.Dataset.Contacts = append(e.Dataset.Contacts, actorToContact(*m.Contact))
	}

	// Editors as associated parties
	for _, actor := range m.Editors {
		e.Dataset.AssociatedParties = append(
			e.Dataset.AssociatedParties,
			actorToAssociatedParty(actor, "editor"),
		)
	}

	// Contributors as associated parties
	for _, actor := range m.Contributors {
		e.Dataset.AssociatedParties = append(
			e.Dataset.AssociatedParties,
			actorToAssociatedParty(actor, "contributor"),
		)
	}

	// Coverage
	if m.GeographicScope != "" || m.TaxonomicScope != "" {
		e.Dataset.Coverage = &dwca.Coverage{}
		if m.GeographicScope != "" {
			e.Dataset.Coverage.GeographicCoverage = &dwca.GeographicCoverage{
				GeographicDescription: m.GeographicScope,
			}
		}
	}

	return e
}

func actorToCreator(a coldp.Actor) dwca.Creator {
	c := dwca.Creator{
		ElectronicMailAddress: a.Email,
	}
	if a.Given != "" || a.Family != "" {
		c.IndividualName = &dwca.IndividualName{
			GivenName: a.Given,
			SurName:   a.Family,
		}
	}
	if a.Organization != "" {
		c.OrganizationName = &dwca.OrganizationName{
			Value: a.Organization,
		}
	}
	return c
}

func actorToContact(a coldp.Actor) dwca.Contact {
	c := dwca.Contact{
		ElectronicMailAddress: a.Email,
	}
	if a.Given != "" || a.Family != "" {
		c.IndividualName = &dwca.IndividualName{
			GivenName: a.Given,
			SurName:   a.Family,
		}
	}
	if a.Organization != "" {
		c.OrganizationName = &dwca.OrganizationName{
			Value: a.Organization,
		}
	}
	if a.Country != "" || a.City != "" {
		c.Address = &dwca.Address{
			Country: a.Country,
			City:    a.City,
		}
	}
	return c
}

func actorToAssociatedParty(a coldp.Actor, role string) dwca.AssociatedParty {
	ap := dwca.AssociatedParty{
		Roles: []string{role},
	}
	if a.Given != "" || a.Family != "" {
		ap.IndividualName = &dwca.IndividualName{
			GivenName: a.Given,
			SurName:   a.Family,
		}
	}
	if a.Organization != "" {
		ap.OrganizationName = &dwca.OrganizationName{
			Value: a.Organization,
		}
	}
	if a.Country != "" || a.City != "" {
		ap.Address = &dwca.Address{
			Country: a.Country,
			City:    a.City,
		}
	}
	return ap
}

func marshalEML(e *dwca.EML) ([]byte, error) {
	lang := e.Lang
	if lang == "" {
		lang = "eng"
	}

	pkgID := e.Dataset.ID
	if pkgID == "" {
		pkgID = "synthetic_" + time.Now().Format("2006-01-02")
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(emlEnvelopeOpen, pkgID, lang))

	enc := xml.NewEncoder(&buf)
	enc.Indent("  ", "  ")

	start := xml.StartElement{
		Name: xml.Name{Local: "dataset"},
	}
	if e.Dataset.ID != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "id"},
			Value: e.Dataset.ID,
		})
	}

	if err := enc.EncodeElement(e.Dataset, start); err != nil {
		return nil, err
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}

	buf.WriteString("\n")
	buf.WriteString(emlEnvelopeClose)

	return buf.Bytes(), nil
}
