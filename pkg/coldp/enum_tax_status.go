package coldp

import "strings"

type TaxonomicStatus int

const (
	UnknownTaxSt TaxonomicStatus = iota
	AcceptedTS
	ProvisionallyAcceptedTS
	SynonymTS
	AmbiguousSynonymTS
	MisappliedTS
	BareNameTS
)

var taxStatToString = map[TaxonomicStatus]string{
	AcceptedTS:              "ACCEPTED",
	ProvisionallyAcceptedTS: "PROVISIONALLY_ACCEPTED",
	SynonymTS:               "SYNONYM",
	AmbiguousSynonymTS:      "AMBIGUOUS_SYNONYM",
	MisappliedTS:            "MISAPPLIED",
	BareNameTS:              "BARE_NAME",
}

var stringToTaxStat = func() map[string]TaxonomicStatus {
	res := make(map[string]TaxonomicStatus)
	for k, v := range taxStatToString {
		res[v] = k
	}
	return res
}()

func (ts TaxonomicStatus) ID() string {
	if res, ok := taxStatToString[ts]; ok {
		return res
	}
	return ""
}

func (ts TaxonomicStatus) String() string {
	return ToStr(ts.ID())
}

func NewTaxonomicStatus(s string) TaxonomicStatus {
	s = strings.ToUpper(s)
	s = strings.Replace(s, " ", "_", -1)
	if res, ok := stringToTaxStat[s]; ok {
		return res
	}
	return UnknownTaxSt
}
