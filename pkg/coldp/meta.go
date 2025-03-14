package coldp

// Meta represents metadata for a dataset.
type Meta struct {
	Key               string   `yaml:"key,omitempty"               json:"key,omitempty"`
	Title             string   `yaml:"title"                       json:"title"`
	Alias             string   `yaml:"alias,omitempty"             json:"alias,omitempty"`
	Description       string   `yaml:"description,omitempty"       json:"description,omitempty"`
	DOI               string   `yaml:"doi,omitempty"               json:"doi,omitempty"`
	Issued            string   `yaml:"issued,omitempty"            json:"issued,omitempty"`
	Version           string   `yaml:"version,omitempty"           json:"version,omitempty"`
	GeographicScope   string   `yaml:"geographicScope,omitempty"   json:"geographicScope,omitempty"`
	TaxonomicScope    string   `yaml:"taxonomicScope,omitempty"    json:"taxonomicScope,omitempty"`
	TemporalScope     string   `yaml:"temporalScope,omitempty"     json:"temporalScope,omitempty"`
	Keywords          []string `yaml:"keywords,omitempty"          json:"keywords,omitempty"`
	Confidence        *int     `yaml:"confidence,omitempty"        json:"confidence,omitempty"`
	Completeness      *int     `yaml:"completeness,omitempty"      json:"completeness,omitempty"`
	License           string   `yaml:"license,omitempty"           json:"license,omitempty"`
	URL               string   `yaml:"url,omitempty"               json:"url,omitempty"`
	Logo              string   `yaml:"logo,omitempty"              json:"logo,omitempty"`
	Label             string   `yaml:"label,omitempty"             json:"label,omitempty"`
	Citation          string   `yaml:"citation,omitempty"          json:"citation,omitempty"`
	LastImportAttempt string   `yaml:"lastImportAttempt,omitempty" json:"lastImportAttempt,omitempty"`
	LastImportState   string   `yaml:"lastImportState,omitempty"   json:"lastImportState,omitempty"`
	Private           bool     `yaml:"private,omitempty"           json:"private,omitempty"`
	Contact           *Actor   `yaml:"contact,omitempty"           json:"contact,omitempty"`
	Publisher         *Actor   `yaml:"publisher,omitempty"         json:"publisher,omitempty"`
	Editors           []Actor  `yaml:"editor,omitempty"            json:"editor,omitempty"`
	Creators          []Actor  `yaml:"creator,omitempty"           json:"creator,omitempty"`
	Contributors      []Actor  `yaml:"contributor,omitempty"       json:"contributor,omitempty"`
	Sources           []Source `yaml:"source,omitempty"            json:"source,omitempty"`
}

// Actor represents an individual or organization.
type Actor struct {
	Orcid        string `yaml:"orcid,omitempty"        json:"orcid,omitempty"`
	Given        string `yaml:"given,omitempty"        json:"given,omitempty"`
	Family       string `yaml:"family,omitempty"       json:"family,omitempty"`
	Email        string `yaml:"email,omitempty"        json:"email,omitempty"`
	City         string `yaml:"city,omitempty"         json:"city,omitempty"`
	State        string `yaml:"state,omitempty"        json:"state,omitempty"`
	Country      string `yaml:"country,omitempty"      json:"country,omitempty"`
	RorID        string `yaml:"rorid,omitempty"        json:"rorid,omitempty"`
	Organization string `yaml:"organisation,omitempty" json:"organisation,omitempty"`
	URL          string `yaml:"url,omitempty"          json:"url,omitempty"`
	Note         string `yaml:"note,omitempty"         json:"note,omitempty"`
}

// Source represents a source of data.
type Source struct {
	ID      string `yaml:"id"               json:"id"`
	Type    string `yaml:"type,omitempty"   json:"type,omitempty"`
	Title   string `yaml:"title,omitempty"  json:"title,omitempty"`
	Authors []any  `yaml:"author,omitempty" json:"author,omitempty"`
	Issued  string `yaml:"issued,omitempty" json:"issued,omitempty"`
	ISBN    string `yaml:"isbn,omitempty"   json:"isbn,omitempty"`
}
