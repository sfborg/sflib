package coldp

import "strings"

type GeoTime int

const (
	UnknownGT GeoTime = iota
	Hadean
	Precambrian
	Archean
	Eoarchean
	Paleoarchean
	Mesoarchean
	Neoarchean
	Proterozoic
	Paleoproterozoic
	Siderian
	Rhyacian
	Orosirian
	Statherian
	Mesoproterozoic
	Calymmian
	Ectasian
	Stenian
	Tonian
	Neoproterozoic
	Cryogenian
	Ediacaran
	Cambrian
	Fortunian
	Paleozoic
	Phanerozoic
	Terreneuvian
	CambrianStage1
	CambrianSeries
	CambrianStage2
	CambrianStage3
	Wuliuan
	Miaolingian
	Drumian
	Guzhangian
	Furongian
	Paibian
	Jiangshanian
	CambrianStage4
	Tremadocian
	LowerOrdovician
	Ordovician
	Floian
	Dapingian
	MiddleOrdovician
	Darriwilian
	Sandbian
	UpperOrdovician
	Katian
	Hirnantian
	Llandovery
	Rhuddanian
	Silurian
	Aeronian
	Telychian
	Sheinwoodian
	Wenlock
	Homerian
	Ludlow
	Gorstian
	Ludfordian
	Pridoli
	Devonian
	LowerDevonian
	Lochkovian
	Pragian
	Emsian
	Eifelian
	MiddleDevonian
	Givetian
	UpperDevonian
	Frasnian
	Famennian
	LowerMississippian
	Tournaisian
	Mississippian
	Carboniferous
	MiddleMississippian
	Visean
	Serpukhovian
	UpperMississippian
	Bashkirian
	Pennsylvanian
	LowerPennsylvanian
	MiddlePennsylvanian
	Moscovian
	Kasimovian
	UpperPennsylvanian
	Gzhelian
	Cisuralian
	Asselian
	Permian
	Sakmarian
	Artinskian
	Kungurian
	Roadian
	Guadalupian
	Wordian
	Capitanian
	Lopingian
	Wuchiapingian
	Changhsingian
	Induan
	LowerTriassic
	Mesozoic
	Triassic
	Olenekian
	Anisian
	MiddleTriassic
	Ladinian
	Carnian
	UpperTriassic
	Norian
	Rhaetian
	Jurassic
	Hettangian
	LowerJurassic
	Sinemurian
	Pliensbachian
	Toarcian
	MiddleJurassic
	Aalenian
	Bajocian
	Bathonian
	Callovian
	Oxfordian
	UpperJurassic
	Kimmeridgian
	Tithonian
	LowerCretaceous
	Cretaceous
	Berriasian
	Valanginian
	Hauterivian
	Barremian
	Aptian
	Albian
	Cenomanian
	UpperCretaceous
	Turonian
	Coniacian
	Santonian
	Campanian
	Maastrichtian
	Paleocene
	Paleogene
	Cenozoic
	Danian
	Selandian
	Thanetian
	Eocene
	Ypresian
	Lutetian
	Bartonian
	Priabonian
	Rupelian
	Oligocene
	Chattian
	Aquitanian
	Neogene
	Miocene
	Burdigalian
	Langhian
	Serravallian
	Tortonian
	Messinian
	Zanclean
	Pliocene
	Piacenzian
	Quaternary
	Gelasian
	Pleistocene
	Calabrian
	MiddlePleistocene
	UpperPleistocene
	Holocene
	Greenlandian
	Northgrippian
	Meghalayan
)

var geotimeToString = map[GeoTime]string{
	Hadean:              "HADEAN",
	Precambrian:         "PRECAMBRIAN",
	Archean:             "ARCHEAN",
	Eoarchean:           "EOARCHEAN",
	Paleoarchean:        "PALEOARCHEAN",
	Mesoarchean:         "MESOARCHEAN",
	Neoarchean:          "NEOARCHEAN",
	Proterozoic:         "PROTEROZOIC",
	Paleoproterozoic:    "PALEOPROTEROZOIC",
	Siderian:            "SIDERIAN",
	Rhyacian:            "RHYACIAN",
	Orosirian:           "OROSIRIAN",
	Statherian:          "STATHERIAN",
	Mesoproterozoic:     "MESOPROTEROZOIC",
	Calymmian:           "CALYMMIAN",
	Ectasian:            "ECTASIAN",
	Stenian:             "STENIAN",
	Tonian:              "TONIAN",
	Neoproterozoic:      "NEOPROTEROZOIC",
	Cryogenian:          "CRYOGENIAN",
	Ediacaran:           "EDIACARAN",
	Cambrian:            "CAMBRIAN",
	Fortunian:           "FORTUNIAN",
	Paleozoic:           "PALEOZOIC",
	Phanerozoic:         "PHANEROZOIC",
	Terreneuvian:        "TERRENEUVIAN",
	CambrianStage1:      "CAMBRIAN_STAGE_1",
	CambrianSeries:      "CAMBRIAN_SERIE_S",
	CambrianStage2:      "CAMBRIAN_STAGE_2",
	CambrianStage3:      "CAMBRIAN_STAGE_3",
	Wuliuan:             "WULIUAN",
	Miaolingian:         "MIAOLINGIAN",
	Drumian:             "DRUMIAN",
	Guzhangian:          "GUZHANGIAN",
	Furongian:           "FURONGIAN",
	Paibian:             "PAIBIAN",
	Jiangshanian:        "JIANGSHANIAN",
	CambrianStage4:      "CAMBRIAN_STAGE_4",
	Tremadocian:         "TREMADOCIAN",
	LowerOrdovician:     "LOWER_ORDOVICIAN",
	Ordovician:          "ORDOVICIAN",
	Floian:              "FLOIAN",
	Dapingian:           "DAPINGIAN",
	MiddleOrdovician:    "MIDDLE_ORDOVICIAN",
	Darriwilian:         "DARRIWILIAN",
	Sandbian:            "SANDBIAN",
	UpperOrdovician:     "UPPER_ORDOVICIAN",
	Katian:              "KATIAN",
	Hirnantian:          "HIRNANTIAN",
	Llandovery:          "LLANDOVERY",
	Rhuddanian:          "RHUDDANIAN",
	Silurian:            "SILURIAN",
	Aeronian:            "AERONIAN",
	Telychian:           "TELYCHIAN",
	Sheinwoodian:        "SHEINWOODIAN",
	Wenlock:             "WENLOCK",
	Homerian:            "HOMERIAN",
	Ludlow:              "LUDLOW",
	Gorstian:            "GORSTIAN",
	Ludfordian:          "LUDFORDIAN",
	Pridoli:             "PRIDOLI",
	Devonian:            "DEVONIAN",
	LowerDevonian:       "LOWER_DEVONIAN",
	Lochkovian:          "LOCHKOVIAN",
	Pragian:             "PRAGIAN",
	Emsian:              "EMSIAN",
	Eifelian:            "EIFELIAN",
	MiddleDevonian:      "MIDDLE_DEVONIAN",
	Givetian:            "GIVETIAN",
	UpperDevonian:       "UPPER_DEVONIAN",
	Frasnian:            "FRASNIAN",
	Famennian:           "FAMENNIAN",
	LowerMississippian:  "LOWER_MISSISSIPPIAN",
	Tournaisian:         "TOURNAISIAN",
	Mississippian:       "MISSISSIPPIAN",
	Carboniferous:       "CARBONIFEROUS",
	MiddleMississippian: "MIDDLE_MISSISSIPPIAN",
	Visean:              "VISEAN",
	Serpukhovian:        "SERPUKHOVIAN",
	UpperMississippian:  "UPPER_MISSISSIPPIAN",
	Bashkirian:          "BASHKIRIAN",
	Pennsylvanian:       "PENNSYLVANIAN",
	LowerPennsylvanian:  "LOWER_PENNSYLVANIAN",
	MiddlePennsylvanian: "MIDDLE_PENNSYLVANIAN",
	Moscovian:           "MOSCOVIAN",
	Kasimovian:          "KASIMOVIAN",
	UpperPennsylvanian:  "UPPER_PENNSYLVANIAN",
	Gzhelian:            "GZHELIAN",
	Cisuralian:          "CISURALIAN",
	Asselian:            "ASSELIAN",
	Permian:             "PERMIAN",
	Sakmarian:           "SAKMARIAN",
	Artinskian:          "ARTINSKIAN",
	Kungurian:           "KUNGURIAN",
	Roadian:             "ROADIAN",
	Guadalupian:         "GUADALUPIAN",
	Wordian:             "WORDIAN",
	Capitanian:          "CAPITANIAN",
	Lopingian:           "LOPINGIAN",
	Wuchiapingian:       "WUCHIAPINGIAN",
	Changhsingian:       "CHANGHSINGIAN",
	Induan:              "INDUAN",
	LowerTriassic:       "LOWER_TRIASSIC",
	Mesozoic:            "MESOZOIC",
	Triassic:            "TRIASSIC",
	Olenekian:           "OLENEKIAN",
	Anisian:             "ANISIAN",
	MiddleTriassic:      "MIDDLE_TRIASSIC",
	Ladinian:            "LADINIAN",
	Carnian:             "CARNIAN",
	UpperTriassic:       "UPPER_TRIASSIC",
	Norian:              "NORIAN",
	Rhaetian:            "RHAETIAN",
	Jurassic:            "JURASSIC",
	Hettangian:          "HETTANGIAN",
	LowerJurassic:       "LOWER_JURASSIC",
	Sinemurian:          "SINEMURIAN",
	Pliensbachian:       "PLIENSBACHIAN",
	Toarcian:            "TOARCIAN",
	MiddleJurassic:      "MIDDLE_JURASSIC",
	Aalenian:            "AALENIAN",
	Bajocian:            "BAJOCIAN",
	Bathonian:           "BATHONIAN",
	Callovian:           "CALLOVIAN",
	Oxfordian:           "OXFORDIAN",
	UpperJurassic:       "UPPER_JURASSIC",
	Kimmeridgian:        "KIMMERIDGIAN",
	Tithonian:           "TITHONIAN",
	LowerCretaceous:     "LOWER_CRETACEOUS",
	Cretaceous:          "CRETACEOUS",
	Berriasian:          "BERRIASIAN",
	Valanginian:         "VALANGINIAN",
	Hauterivian:         "HAUTERIVIAN",
	Barremian:           "BARREMIAN",
	Aptian:              "APTIAN",
	Albian:              "ALBIAN",
	Cenomanian:          "CENOMANIAN",
	UpperCretaceous:     "UPPER_CRETACEOUS",
	Turonian:            "TURONIAN",
	Coniacian:           "CONIACIAN",
	Santonian:           "SANTONIAN",
	Campanian:           "CAMPANIAN",
	Maastrichtian:       "MAASTRICHTIAN",
	Paleocene:           "PALEOCENE",
	Paleogene:           "PALEOGENE",
	Cenozoic:            "CENOZOIC",
	Danian:              "DANIAN",
	Selandian:           "SELANDIAN",
	Thanetian:           "THANETIAN",
	Eocene:              "EOCENE",
	Ypresian:            "YPRESIAN",
	Lutetian:            "LUTETIAN",
	Bartonian:           "BARTONIAN",
	Priabonian:          "PRIABONIAN",
	Rupelian:            "RUPELIAN",
	Oligocene:           "OLIGOCENE",
	Chattian:            "CHATTIAN",
	Aquitanian:          "AQUITANIAN",
	Neogene:             "NEOGENE",
	Miocene:             "MIOCENE",
	Burdigalian:         "BURDIGALIAN",
	Langhian:            "LANGHIAN",
	Serravallian:        "SERRAVALLIAN",
	Tortonian:           "TORTONIAN",
	Messinian:           "MESSINIAN",
	Zanclean:            "ZANCLEAN",
	Pliocene:            "PLIOCENE",
	Piacenzian:          "PIACENZIAN",
	Quaternary:          "QUATERNARY",
	Gelasian:            "GELASIAN",
	Pleistocene:         "PLEISTOCENE",
	Calabrian:           "CALABRIAN",
	MiddlePleistocene:   "MIDDLE_PLEISTOCENE",
	UpperPleistocene:    "UPPERP_LEISTOCENE",
	Holocene:            "HOLOCENE",
	Greenlandian:        "GREENLANDIAN",
	Northgrippian:       "NORTHGRIPPIAN",
	Meghalayan:          "MEGHALAYAN",
}

var stringToGeoTime = func() map[string]GeoTime {
	res := make(map[string]GeoTime)
	for k, v := range geotimeToString {
		res[v] = k
	}
	return res
}()

// ID return the string ID of GeoTime.
func (g GeoTime) ID() string {
	if res, ok := geotimeToString[g]; ok {
		return res
	}
	return ""
}

// String return the string representation of GeoTime.
func (g GeoTime) String() string {
	return ToTitleCase(g.ID())
}

func NewGeoTime(s string) GeoTime {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, " ", "_")
	if res, ok := stringToGeoTime[s]; ok {
		return res
	}
	return UnknownGT
}
