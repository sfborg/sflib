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
		return "full"
	case SciNameCanonical:
		return "canonical"
	case SciNameComposite:
		return "composite"
	default:
		return "unknown"
	}
}

type SynonymType int

const (
	SynUnknown SynonymType = iota
	SynAcceptedID
	SynHierarchy
	SynExtension
)

func (st SynonymType) String() string {
	switch st {
	case SynAcceptedID:
		return "accepted ID"
	case SynHierarchy:
		return "hierarchy"
	case SynExtension:
		return "extension"
	default:
		return "unknown"
	}
}

type HierType int

const (
	HierUnknown HierType = iota
	HierTree
	HierFlat
)

func (h HierType) String() string {
	switch h {
	case HierTree:
		return "tree"
	case HierFlat:
		return "flat"
	default:
		return "unknown"
	}
}
