package coldp

import (
	"strings"
)

type CSLRecords []CSL

type CSL struct {
	ID                  string       `json:"id"`
	Publisher           string       `json:"publisher,omitempty"`
	Issue               string       `json:"issue,omitempty"`
	PublishedPrint      Date         `json:"published-print,omitempty"`
	DOI                 string       `json:"DOI,omitempty"`
	Type                string       `json:"type,omitempty"`
	Created             Date         `json:"created,omitempty"`
	Page                string       `json:"page,omitempty"`
	Source              string       `json:"source,omitempty"`
	Title               string       `json:"title,omitempty"`
	Prefix              string       `json:"prefix,omitempty"`
	Volume              string       `json:"volume,omitempty"`
	Authors             []AuthorCSL  `json:"author,omitempty"`
	ContainerTitle      string       `json:"container-title,omitempty"`
	ContainerTitleShort any          `json:"container-title-short,omitempty"`
	OriginalTitle       any          `json:"original-title,omitempty"`
	Language            string       `json:"language,omitempty"`
	Links               []Link       `json:"link,omitempty"`
	Deposited           Date         `json:"deposited,omitempty"`
	Subtitle            any          `json:"subtitle,omitempty"`
	ShortTitle          any          `json:"short-title,omitempty"`
	Issued              Date         `json:"issued,omitempty"`
	JournalIssue        JournalIssue `json:"journal-issue,omitempty"`
	URL                 string       `json:"URL,omitempty"`
	ISSN                any          `json:"ISSN,omitempty"`
	Subject             any          `json:"subject,omitempty"`
}

type Date struct {
	DateParts []DatePart `json:"date-parts,omitempty"`
	DateTime  string     `json:"date-time,omitempty"`
	Timestamp string     `json:"timestamp,omitempty"`
}

type DatePart struct {
	Date []string
}

type AuthorCSL struct {
	Given        string        `json:"given,omitempty"`
	Family       string        `json:"family,omitempty"`
	Sequence     string        `json:"sequence,omitempty"`
	Affiliations []Affiliation `json:"affiliation,omitempty"`
}

type Affiliation struct {
}

type Link struct {
	URL                 string `json:"URL,omitempty"`
	ContentType         string `json:"content-type,omitempty"`
	ContentVersion      string `json:"content-version,omitempty"`
	IntendedApplication string `json:"similarity-checking,omitempty"`
}

type JournalIssue struct {
	PublishedPrint Date   `json:"published-print,omitempty"`
	Issue          string `json:"issue,omitempty"`
}

func (c CSL) ToReference() DataLoader {
	res := Reference{
		ID:                  c.ID,
		Type:                NewReferenceType(c.Type),
		Author:              flattenAuthorsCSL(c.Authors),
		Title:               c.Title,
		TitleShort:          fromAny(c.ShortTitle),
		ContainerTitle:      c.ContainerTitle,
		ContainerTitleShort: fromAny(c.ContainerTitleShort),
		Issued:              fromDate(c.Issued),
		Volume:              c.Volume,
		Issue:               c.Issue,
		Page:                c.Page,
		Publisher:           c.Publisher,
		ISSN:                fromAny(c.ISSN),
		DOI:                 c.DOI,
		Link:                c.URL,
	}
	return res
}

func flattenAuthorsCSL(authors []AuthorCSL) string {
	res := []string{}
	for _, a := range authors {
		name := a.Family
		if a.Given != "" {
			name += "," + a.Given
		}
		res = append(res, name)
	}
	return strings.Join(res, ", ")
}

func fromAny(d any) string {
	switch v := d.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ",")
	default:
		return ""
	}
}

func fromDate(d Date) string {
	if d.DateTime != "" {
		return d.DateTime
	}

	if len(d.DateParts) > 0 {
		return strings.Join(d.DateParts[0].Date, "-")
	}

	return ""
}
