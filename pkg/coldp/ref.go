package coldp

type Reference struct {
	// ID of the reference.
	ID string

	// AlternativeID is a list of ids separated by comma (scope:id, id).
	AlternativeID string

	// SourceID to source from metadata.
	SourceID string

	// Citation unparsed string for reference's citation.
	Citation string

	// Type of the publication.
	Type ReferenceType

	// Author is the list of author/s.
	Author string

	// AuthorID is the comma-separated list of author ids.
	AuthorID string

	// Editor is the list of editor/s.
	Editor string

	// EditorID is the comma-separated list of editorIDs.
	EditorID string

	// Title of the reference.
	Title string

	// TitleShort of the reference.
	TitleShort string

	// ContainerAuthor is author(s) of the container holding the item, e.g. the
	// book author for a book chapter. See author for recommendations how to
	// supply person names.
	ContainerAuthor string

	// ContainerAuthorID is the comma-separated list of author/s of the
	// container.
	ContainerAuthorID string

	// ContainerTitle is the title of the container.
	ContainerTitle string

	// ContainerTitleShort of the container.
	ContainerTitleShort string

	// Issued is ISO8601 date Date the work was issued/published. Date
	// can be truncated if month or day are absent.
	Issued string

	// Accessed is ISO8601 date Date the work was visited. Date
	// can be truncated if month or day are absent.
	Accessed string

	// CollectionTitle is the  title of the collection holding the item, e.g. the
	// series title for a book.
	CollectionTitle string

	// CollectionEditor is the editor of the collection.
	CollectionEditor string

	// Volume of the reference.
	Volume string

	// Issue of the reference.
	Issue string

	// Edition of the reference.
	Edition string

	// Page of the reference.
	Page string

	// Publisher of the refernce.
	Publisher string

	// PublisherPlace is the geographic location of the publisher.
	PublisherPlace string

	// Version of the reference (if applicable).
	Version string

	// ISBN is the International Standard Book Number
	ISBN string

	// ISSN is the International Standard Serial Number
	ISSN string

	// DOI of the reference
	DOI string

	// Link is the URL to the reference.
	Link string

	// Remarks about the reference.
	Remarks string

	// Modified is a timestamp.
	Modified string

	// ModifiedBy a person who last modified the record.
	ModifiedBy string
}

func (r Reference) Headers() []string {
	return []string{
		"col:id",
		"col:alternativeId",
		"col:sourceId",
		"col:citation",
		"col:type",
		"col:author",
		"col:authorId",
		"col:editor",
		"col:editorId",
		"col:title",
		"col:titleShort",
		"col:containerAuthor",
		"col:containerTitle",
		"col:containerTitleShort",
		"col:issued",
		"col:accessed",
		"col:collectionTitle",
		"col:collectionEditor",
		"col:volume",
		"col:issue",
		"col:edition",
		"col:page",
		"col:publisher",
		"col:publisherPlace",
		"col:version",
		"col:isbn",
		"col:issn",
		"col:doi",
		"col:link",
		"col:remarks",
		"col:modified",
		"col:modifiedBy",
	}
}

func (r Reference) Load(headers, data []string) (DataLoader, error) {
	row, warning := RowToMap(headers, data)
	r.ID = row["id"]
	r.AlternativeID = row["alternativeid"]
	r.SourceID = row["sourceid"]
	r.Citation = row["citation"]
	r.Type = NewReferenceType(row["type"])
	r.Author = row["author"]
	r.AuthorID = row["authorid"]
	r.Editor = row["editor"]
	r.EditorID = row["editorid"]
	r.Title = row["title"]
	r.TitleShort = row["titleshort"]
	r.ContainerAuthor = row["containerauthor"]
	r.ContainerTitle = row["containertitle"]
	r.ContainerTitleShort = row["containertitleshort"]
	r.Issued = row["issued"]
	r.Accessed = row["accessed"]
	r.CollectionTitle = row["collectiontitle"]
	r.CollectionEditor = row["collectioneditor"]
	r.Volume = row["volume"]
	r.Issue = row["issue"]
	r.Edition = row["edition"]
	r.Page = row["page"]
	r.Publisher = row["publisher"]
	r.PublisherPlace = row["publisherplace"]
	r.Version = row["version"]
	r.ISBN = row["isbn"]
	r.ISSN = row["issn"]
	r.DOI = row["doi"]
	r.Link = row["link"]
	r.Remarks = row["remarks"]
	r.Modified = row["modified"]
	r.ModifiedBy = row["modifiedby"]
	return r, warning
}
