package dwca

import "encoding/xml"

type EML struct {
	XMLName        xml.Name `xml:"eml"`
	Lang           string   `xml:"xml:lang,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Dataset        Dataset  `xml:"dataset"`
}

type Dataset struct {
	ID                    string              `xml:"id,attr"`
	AlternativeIdentifier AltID               `xml:"alternateIdentifier"`
	Title                 string              `xml:"title"`
	Creators              []Creator           `xml:"creator"`
	MetadataProviders     []MetadataProvider  `xml:"metadataProvider"`
	AssociatedParties     []AssociatedParty   `xml:"associatedParty"`
	PubDate               string              `xml:"pubDate"`
	Language              string              `xml:"language"`
	Abstract              Abstract            `xml:"abstract"`
	KeywodSets            []KeywordSet        `xml:"keywordSet"`
	IntellectualRights    *IntellectualRights `xml:"intellectualRights"`
	Coverage              *Coverage           `xml:"coverage"`
	Contacts              []Contact           `xml:"contact"`
}

type Contact struct {
	IndividualName        *IndividualName   `xml:"individualName"`
	OrganizationName      *OrganizationName `xml:"organizationName"`
	Address               *Address          `xml:"address,omitempty"`
	ElectronicMailAddress string            `xml:"electronicMailAddress"`
}

type Coverage struct {
	GeographicCoverage *GeographicCoverage `xml:"geographicCoverage"`
	TemporalCoverage   *TemporalCoverage   `xml:"temporalCoverage"`
}

type TemporalCoverage struct {
	BeginDate CalendarDate `xml:"beginDate"`
	EndDate   CalendarDate `xml:"endDate"`
}

type CalendarDate struct {
	Value string `xml:",chardata"`
}

type GeographicCoverage struct {
	GeographicDescription string               `xml:"geographicDescription,omitempty"`
	BoundingCoordinates   *BoundingCoordinates `xml:"boundingCoordinates,omitempty"`
}

type BoundingCoordinates struct {
	WestBoundingCoordinate  float64 `xml:"westBoundingCoordinate"`
	EastBoundingCoordinate  float64 `xml:"eastBoundingCoordinate"`
	SouthBoundingCoordinate float64 `xml:"southBoundingCoordinate"`
	NorthBoundingCoordinate float64 `xml:"northBoundingCoordinate"`
}

type IntellectualRights struct {
	Para string `xml:"para"`
}

type KeywordSet struct {
	Keywords []Keyword `xml:"keyword"`
}

type Keyword struct {
	Value string `xml:",chardata"`
}

type AssociatedParty struct {
	IndividualName   *IndividualName   `xml:"individualName"`
	OrganizationName *OrganizationName `xml:"organizationName"`
	Address          *Address          `xml:"address,omitempty"`
	Roles            []string          `xml:"role"`
}

type AltID struct {
	XMLName xml.Name `xml:"alternateIdentifier"`
	System  string   `xml:"system,attr,omitempty"`
	Value   string   `xml:",chardata"`
}

type Creator struct {
	ID                    string `xml:"id,attr"`
	Scope                 string `xml:"scope,attr,omitempty"`
	IndividualName        *IndividualName
	OrganizationName      *OrganizationName `xml:"organizationName"`
	ElectronicMailAddress string            `xml:"electronicMailAddress"`
}

type OrganizationName struct {
	XMLName xml.Name `xml:"organizationName"`
	Value   string   `xml:",chardata"`
}

type MetadataProvider struct {
	IndividualName        *IndividualName   `xml:"individualName"`
	OrganizationName      *OrganizationName `xml:"organizationName"`
	Address               *Address          `xml:"address,omitempty"`
	ElectronicMailAddress string            `xml:"electronicMailAddress"`
	OnlineURL             string            `xml:"onlineUrl,omitempty"`
}

type Address struct {
	Country    string `xml:"country,omitempty"`
	City       string `xml:"city,omitempty"`
	PostalCode string `xml:"postalCode,omitempty"`
}

type IndividualName struct {
	XMLName   xml.Name `xml:"individualName"`
	GivenName string   `xml:"givenName"`
	SurName   string   `xml:"surName"`
}

type Abstract struct {
	Para string `xml:"para"`
}
