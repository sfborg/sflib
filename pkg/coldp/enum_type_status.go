package coldp

import "strings"

type TypeStatus int

const (
	UnknownTS TypeStatus = iota
	OtherTS
	HomoeotypeTS
	PlesiotypeTS
	PlastotypeTS
	PlastosyntypeTS
	PlastoparatypeTS
	PlastoneotypeTS
	PlastolectotypeTS
	PlastoisotypeTS
	PlastoholotypeTS
	AllotypeTS
	AlloneotypeTS
	AllolectotypeTS
	ParaneotypeTS
	ParalectotypeTS
	IsosyntypeTS
	IsoparatypeTS
	IsoneotypeTS
	IsolectotypeTS
	IsoepitypeTS
	IsotypeTS
	TopotypeTS
	SyntypeTS
	PathotypeTS
	ParatypeTS
	OriginalTS
	NeotypeTS
	LectotypeTS
	IconotypeTS
	HolotypeTS
	HapantotypeTS
	ExTS
	ErgatotypeTS
	EpitypeTS
)

var typeStatusToString = map[TypeStatus]string{
	OtherTS:           "OTHER",
	HomoeotypeTS:      "HOMOEOTYPE",
	PlesiotypeTS:      "PLESIOTYPE",
	PlastotypeTS:      "PLASTOTYPE",
	PlastosyntypeTS:   "PLASTOSYNTYPE",
	PlastoparatypeTS:  "PLASTOPARATYPE",
	PlastoneotypeTS:   "PLASTONEOTYPE",
	PlastolectotypeTS: "PLASTOLECTOTYPE",
	PlastoisotypeTS:   "PLASTOISOTYPE",
	PlastoholotypeTS:  "PLASTOHOLOTYPE",
	AllotypeTS:        "ALLOTYPE",
	AlloneotypeTS:     "ALLONEOTYPE",
	AllolectotypeTS:   "ALLOLECTOTYPE",
	ParaneotypeTS:     "PARANEOTYPE",
	ParalectotypeTS:   "PARALECTOTYPE",
	IsosyntypeTS:      "ISOSYNTYPE",
	IsoparatypeTS:     "ISOPARATYPE",
	IsoneotypeTS:      "ISONEOTYPE",
	IsolectotypeTS:    "ISOLECTOTYPE",
	IsoepitypeTS:      "ISOEPITYPE",
	IsotypeTS:         "ISOTYPE",
	TopotypeTS:        "TOPOTYPE",
	SyntypeTS:         "SYNTYPE",
	PathotypeTS:       "PATHOTYPE",
	ParatypeTS:        "PARATYPE",
	OriginalTS:        "ORIGINAL",
	NeotypeTS:         "NEOTYPE",
	LectotypeTS:       "LECTOTYPE",
	IconotypeTS:       "ICONOTYPE",
	HolotypeTS:        "HOLOTYPE",
	HapantotypeTS:     "HAPANTOTYPE",
	ExTS:              "EX",
	ErgatotypeTS:      "ERGATOTYPE",
	EpitypeTS:         "EPITYPE",
}

var stringToTypeStatus = func() map[string]TypeStatus {
	res := make(map[string]TypeStatus)
	for k, v := range typeStatusToString {
		res[v] = k
	}
	return res
}()

func (ts TypeStatus) ID() string {
	if res, ok := typeStatusToString[ts]; ok {
		return res
	}
	return ""
}

func (ts TypeStatus) String() string {
	return ToStr(ts.ID())
}

func NewTypeStatus(s string) TypeStatus {
	s = strings.ToUpper(s)
	if res, ok := stringToTypeStatus[s]; ok {
		return res
	}
	return UnknownTS
}
