package coldp

import (
	"log/slog"
	"strings"
)

// NomStatus represents the nomenclatural status of a scientific name.
type NomStatus int

// Constants for different nomenclatural statuses.
const (
	UnknownNomStatus NomStatus = iota
	Established                // The name is validly published and available.
	NotEstablished             // The name is not validly published or unavailable.
	Acceptable                 // The name is potentially valid or acceptable for use.
	Unacceptable               // The name is illegitimate or unacceptable for use.
	Conserved                  // The name is conserved against a competing name.
	Rejected                   // The name is rejected in favor of a competing name.
	Doubtful                   // The application of the name is uncertain.
	Manuscript                 // The name is only published in a manuscript.
	Chresonym                  // The name is a chresonym.
)

// NewNomStatus creates a new NomStatus from a string representation.
// It handles various synonyms and normalizations to ensure consistent matching.
func NewNomStatus(s string) NomStatus {
	sOrig := s
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "_", "")
	switch s {
	case "":
		return UnknownNomStatus
	case "nomenvalidum", "available", "established":
		return Established
	case "nominval", "nomeninvalidum", "unavailable", "notestablished":
		return NotEstablished
	case "nomenlegitimum", "potentiallyvalid", "acceptable":
		return Acceptable
	case "nomilleg", "nomenillegitimum", "objectivelyinvalid", "unacceptable":
		return Unacceptable
	case "nomcons", "nomenconservandum", "conservedname", "conserved":
		return Conserved
	case "nomrej", "nomenrejiciendum", "rejected":
		return Rejected
	case "nomdub", "nomeddubium", "doubtful":
		return Doubtful
	case "manuscript name", "manuscript":
		return Manuscript
	case "chresonym":
		return Chresonym
	default:
		slog.Warn("Cannot find nom. status", "input", sOrig)
		return UnknownNomStatus
	}
}

var nomStatusToString = map[NomStatus]string{
	Established:    "ESTABLISHED",
	NotEstablished: "NOT_ESTABLISHED",
	Acceptable:     "ACCEPTABLE",
	Unacceptable:   "UNACCEPTABLE",
	Conserved:      "CONSERVED",
	Rejected:       "REJECTED",
	Doubtful:       "DOUBTFUL",
	Manuscript:     "MANUSCRIPT",
	Chresonym:      "CHRESONYM",
}

func (n NomStatus) ID() string {
	if res, ok := nomStatusToString[n]; ok {
		return res
	}
	return ""
}

func (n NomStatus) String() string {
	return ToStr(n.ID())
}
