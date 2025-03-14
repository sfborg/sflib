package coldp

type SpInteractionType int

const (
	UnknownSpIntT SpInteractionType = iota
	MutualistOf
	CommensalistOf
	HasEpiphyte
	EpiphyteOf
	HasEggsLayedOnBy
	LaysEggsOn
	PollinatedBy
	Pollinates
	FlowersVisitedBy
	VisitsFlowersOf
	VisitedBy
	Visits
	HasHyperparasitoid
	HyperparasitoidOf
	HasParasitoid
	ParasitoidOf
	HasKleptoparasite
	KleptoparasiteOf
	HasHyperparasite
	HyperparasiteOf
	HasEctoparasite
	EctoparasiteOf
	HasEndoparasite
	EndoparasiteOf
	HasVector
	VectorOf
	HasPathogen
	PathogenOf
	HasParasite
	ParasiteOf
	HasHost
	HostOf
	PreyedUponBy
	PreysUpon
	KilledBy
	Kills
	EatenBy
	Eats
	SymbiontOf
	AdjacentTo
	InteractsWith
	CoOccursWith
	RelatedTo
)

var spIntTypeToString = map[SpInteractionType]string{
	MutualistOf:        "MUTUALIST_OF",
	CommensalistOf:     "COMMENSALIST_OF",
	HasEpiphyte:        "HAS_EPIPHYTE",
	EpiphyteOf:         "EPIPHYTE_OF",
	HasEggsLayedOnBy:   "HAS_EGGS_LAYED_ON_BY",
	LaysEggsOn:         "LAYS_EGGS_ON",
	PollinatedBy:       "POLLINATED_BY",
	Pollinates:         "POLLINATES",
	FlowersVisitedBy:   "FLOWERS_VISITED_BY",
	VisitsFlowersOf:    "VISITS_FLOWERS_OF",
	VisitedBy:          "VISITED_BY",
	Visits:             "VISITS",
	HasHyperparasitoid: "HAS_HYPERPARASITOID",
	HyperparasitoidOf:  "HYPERPARASITOID_OF",
	HasParasitoid:      "HAS_PARASITOID",
	ParasitoidOf:       "PARASITOID_OF",
	HasKleptoparasite:  "HAS_KLEPTOPARASITE",
	KleptoparasiteOf:   "KLEPTOPARASITE_OF",
	HasHyperparasite:   "HAS_HYPERPARASITE",
	HyperparasiteOf:    "HYPERPARASITE_OF",
	HasEctoparasite:    "HAS_ECTOPARASITE",
	EctoparasiteOf:     "ECTOPARASITE_OF",
	HasEndoparasite:    "HAS_ENDOPARASITE",
	EndoparasiteOf:     "ENDOPARASITE_OF",
	HasVector:          "HAS_VECTOR",
	VectorOf:           "VECTOR_OF",
	HasPathogen:        "HAS_PATHOGEN",
	PathogenOf:         "PATHOGEN_OF",
	HasParasite:        "HAS_PARASITE",
	ParasiteOf:         "PARASITE_OF",
	HasHost:            "HAS_HOST",
	HostOf:             "HOST_OF",
	PreyedUponBy:       "PREYED_UPON_BY",
	PreysUpon:          "PREYS_UPON",
	KilledBy:           "KILLED_BY",
	Kills:              "KILLS",
	EatenBy:            "EATEN_BY",
	Eats:               "EATS",
	SymbiontOf:         "SYMBIONT_OF",
	AdjacentTo:         "ADJACENT_TO",
	InteractsWith:      "INTERACTS_WITH",
	CoOccursWith:       "CO_OCCURS_WITH",
	RelatedTo:          "RELATED_TO",
}

var stringToSpIntType = func() map[string]SpInteractionType {
	res := make(map[string]SpInteractionType)
	for k, v := range spIntTypeToString {
		res[v] = k
	}
	return res
}()

func (si SpInteractionType) ID() string {
	if res, ok := spIntTypeToString[si]; ok {
		return res
	}
	return ""
}

func (si SpInteractionType) String() string {
	return ToStr(si.ID())
}

func NewSpInteractionType(s string) SpInteractionType {
	if res, ok := stringToSpIntType[s]; ok {
		return res
	}
	return UnknownSpIntT
}
