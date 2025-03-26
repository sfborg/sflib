package dwca

// MetaSimple is a simplifiec version of Meta object, that is used to access
// metadata fields by their names or indices.
type MetaSimple struct {
	// CoreData is a simplified version of Core data of the DwCA archive.
	CoreData

	// ExtensionsData is a simplified version of Extensions data of the DwCA.
	ExtensionsData map[string]ExtensionData
}

// CoreData is a simplified version of Core data of the DwCA archive.
type CoreData struct {
	// Index is the field index of the Core's Term.
	Index int

	// Term is the field name of the main Core Data (Topic).
	Term string

	// TermFull is the URI of the main Core Data (Topic).
	TermFull string

	// Location is the location of the Core file.
	Location string

	// FieldsData is a map of field Terms to their FieldData.
	FieldsData map[string]FieldData

	// FieldsIdx is a map of field indices to their FieldData.
	FieldsIdx map[int]FieldData
}

// ExtensionData is a simplified version of Extensions data of the DwCA.
type ExtensionData struct {
	// CoreIndex is the index of the Core main field in the Extension.
	// It allows to create a star schema of the DwCA archive.
	CoreIndex int

	// Location is the location of the Extension file.
	Location string

	// FieldsData is a map of field Terms to their FieldData.
	FieldsData map[string]FieldData

	// FieldsIdx is a map of field indices to their FieldData.
	FieldsIdx map[int]FieldData
}

// FieldData is a simplified version of a field in the DwCA archive.
type FieldData struct {
	// Index is the index of the field in the DwCA archive.
	Index int

	// Term is the field name of the field.
	Term string

	// TermFull is the URI of the field.
	TermFull string
}
