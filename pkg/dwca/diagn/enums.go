package diagn

type SciNameType int

const (
	SciNameUnknown SciNameType = iota
	SciNameFull
	SciNameCanonical
	SciNameComposite
)

func (s SciNameType) String() string {
	switch s {
	case SciNameFull:
		return "full name with authorship"
	case SciNameCanonical:
		return "canonical name, no authorship"
	case SciNameComposite:
		return "name split by several fields"
	default:
		return "unknown type of scientific name"
	}
}

type SynonymType int

const (
	SynUnknown SynonymType = iota
	SynAcceptedID
	SynHierarchy
	SynExtension
	SynNone
)

func (st SynonymType) String() string {
	switch st {
	case SynAcceptedID:
		return "synonymy by accepted ID"
	case SynHierarchy:
		return "synonymy by parent ID"
	case SynExtension:
		return "synonymy in an extension file"
	case SynNone:
		return "synonymy not found"
	default:
		return "unknown synonymy"
	}
}

type HierType int

const (
	HierUnknown HierType = iota
	HierTree
	HierFlat
	HierBoth
)

func (h HierType) String() string {
	switch h {
	case HierTree:
		return "tree hierarhcy"
	case HierFlat:
		return "flat hierarchy"
	case HierBoth:
		return "both tree and flat hierarchies"
	default:
		return "unknown"
	}
}
