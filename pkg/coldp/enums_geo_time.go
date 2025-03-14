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
	Hadean:              "Hadean",
	Precambrian:         "Precambrian",
	Archean:             "Archean",
	Eoarchean:           "Eoarchean",
	Paleoarchean:        "Paleoarchean",
	Mesoarchean:         "Mesoarchean",
	Neoarchean:          "Neoarchean",
	Proterozoic:         "Proterozoic",
	Paleoproterozoic:    "Paleoproterozoic",
	Siderian:            "Siderian",
	Rhyacian:            "Rhyacian",
	Orosirian:           "Orosirian",
	Statherian:          "Statherian",
	Mesoproterozoic:     "Mesoproterozoic",
	Calymmian:           "Calymmian",
	Ectasian:            "Ectasian",
	Stenian:             "Stenian",
	Tonian:              "Tonian",
	Neoproterozoic:      "Neoproterozoic",
	Cryogenian:          "Cryogenian",
	Ediacaran:           "Ediacaran",
	Cambrian:            "Cambrian",
	Fortunian:           "Fortunian",
	Paleozoic:           "Paleozoic",
	Phanerozoic:         "Phanerozoic",
	Terreneuvian:        "Terreneuvian",
	CambrianStage1:      "CambrianStage1",
	CambrianSeries:      "CambrianSeries",
	CambrianStage2:      "CambrianStage2",
	CambrianStage3:      "CambrianStage3",
	Wuliuan:             "Wuliuan",
	Miaolingian:         "Miaolingian",
	Drumian:             "Drumian",
	Guzhangian:          "Guzhangian",
	Furongian:           "Furongian",
	Paibian:             "Paibian",
	Jiangshanian:        "Jiangshanian",
	CambrianStage4:      "CambrianStage4",
	Tremadocian:         "Tremadocian",
	LowerOrdovician:     "LowerOrdovician",
	Ordovician:          "Ordovician",
	Floian:              "Floian",
	Dapingian:           "Dapingian",
	MiddleOrdovician:    "MiddleOrdovician",
	Darriwilian:         "Darriwilian",
	Sandbian:            "Sandbian",
	UpperOrdovician:     "UpperOrdovician",
	Katian:              "Katian",
	Hirnantian:          "Hirnantian",
	Llandovery:          "Llandovery",
	Rhuddanian:          "Rhuddanian",
	Silurian:            "Silurian",
	Aeronian:            "Aeronian",
	Telychian:           "Telychian",
	Sheinwoodian:        "Sheinwoodian",
	Wenlock:             "Wenlock",
	Homerian:            "Homerian",
	Ludlow:              "Ludlow",
	Gorstian:            "Gorstian",
	Ludfordian:          "Ludfordian",
	Pridoli:             "Pridoli",
	Devonian:            "Devonian",
	LowerDevonian:       "LowerDevonian",
	Lochkovian:          "Lochkovian",
	Pragian:             "Pragian",
	Emsian:              "Emsian",
	Eifelian:            "Eifelian",
	MiddleDevonian:      "MiddleDevonian",
	Givetian:            "Givetian",
	UpperDevonian:       "UpperDevonian",
	Frasnian:            "Frasnian",
	Famennian:           "Famennian",
	LowerMississippian:  "LowerMississippian",
	Tournaisian:         "Tournaisian",
	Mississippian:       "Mississippian",
	Carboniferous:       "Carboniferous",
	MiddleMississippian: "MiddleMississippian",
	Visean:              "Visean",
	Serpukhovian:        "Serpukhovian",
	UpperMississippian:  "UpperMississippian",
	Bashkirian:          "Bashkirian",
	Pennsylvanian:       "Pennsylvanian",
	LowerPennsylvanian:  "LowerPennsylvanian",
	MiddlePennsylvanian: "MiddlePennsylvanian",
	Moscovian:           "Moscovian",
	Kasimovian:          "Kasimovian",
	UpperPennsylvanian:  "UpperPennsylvanian",
	Gzhelian:            "Gzhelian",
	Cisuralian:          "Cisuralian",
	Asselian:            "Asselian",
	Permian:             "Permian",
	Sakmarian:           "Sakmarian",
	Artinskian:          "Artinskian",
	Kungurian:           "Kungurian",
	Roadian:             "Roadian",
	Guadalupian:         "Guadalupian",
	Wordian:             "Wordian",
	Capitanian:          "Capitanian",
	Lopingian:           "Lopingian",
	Wuchiapingian:       "Wuchiapingian",
	Changhsingian:       "Changhsingian",
	Induan:              "Induan",
	LowerTriassic:       "LowerTriassic",
	Mesozoic:            "Mesozoic",
	Triassic:            "Triassic",
	Olenekian:           "Olenekian",
	Anisian:             "Anisian",
	MiddleTriassic:      "MiddleTriassic",
	Ladinian:            "Ladinian",
	Carnian:             "Carnian",
	UpperTriassic:       "UpperTriassic",
	Norian:              "Norian",
	Rhaetian:            "Rhaetian",
	Jurassic:            "Jurassic",
	Hettangian:          "Hettangian",
	LowerJurassic:       "LowerJurassic",
	Sinemurian:          "Sinemurian",
	Pliensbachian:       "Pliensbachian",
	Toarcian:            "Toarcian",
	MiddleJurassic:      "MiddleJurassic",
	Aalenian:            "Aalenian",
	Bajocian:            "Bajocian",
	Bathonian:           "Bathonian",
	Callovian:           "Callovian",
	Oxfordian:           "Oxfordian",
	UpperJurassic:       "UpperJurassic",
	Kimmeridgian:        "Kimmeridgian",
	Tithonian:           "Tithonian",
	LowerCretaceous:     "LowerCretaceous",
	Cretaceous:          "Cretaceous",
	Berriasian:          "Berriasian",
	Valanginian:         "Valanginian",
	Hauterivian:         "Hauterivian",
	Barremian:           "Barremian",
	Aptian:              "Aptian",
	Albian:              "Albian",
	Cenomanian:          "Cenomanian",
	UpperCretaceous:     "UpperCretaceous",
	Turonian:            "Turonian",
	Coniacian:           "Coniacian",
	Santonian:           "Santonian",
	Campanian:           "Campanian",
	Maastrichtian:       "Maastrichtian",
	Paleocene:           "Paleocene",
	Paleogene:           "Paleogene",
	Cenozoic:            "Cenozoic",
	Danian:              "Danian",
	Selandian:           "Selandian",
	Thanetian:           "Thanetian",
	Eocene:              "Eocene",
	Ypresian:            "Ypresian",
	Lutetian:            "Lutetian",
	Bartonian:           "Bartonian",
	Priabonian:          "Priabonian",
	Rupelian:            "Rupelian",
	Oligocene:           "Oligocene",
	Chattian:            "Chattian",
	Aquitanian:          "Aquitanian",
	Neogene:             "Neogene",
	Miocene:             "Miocene",
	Burdigalian:         "Burdigalian",
	Langhian:            "Langhian",
	Serravallian:        "Serravallian",
	Tortonian:           "Tortonian",
	Messinian:           "Messinian",
	Zanclean:            "Zanclean",
	Pliocene:            "Pliocene",
	Piacenzian:          "Piacenzian",
	Quaternary:          "Quaternary",
	Gelasian:            "Gelasian",
	Pleistocene:         "Pleistocene",
	Calabrian:           "Calabrian",
	MiddlePleistocene:   "MiddlePleistocene",
	UpperPleistocene:    "UpperPleistocene",
	Holocene:            "Holocene",
	Greenlandian:        "Greenlandian",
	Northgrippian:       "Northgrippian",
	Meghalayan:          "Meghalayan",
}

var stringToGeoTime = func() map[string]GeoTime {
	res := make(map[string]GeoTime)
	for k, v := range geotimeToString {
		v = strings.ToLower(v)
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
	return ToStr(g.ID())
}

func NewGeoTime(s string) GeoTime {
	s = strings.ToLower(s)
	s = strings.ReplaceAll("s", " ", "")
	if res, ok := stringToGeoTime[s]; ok {
		return res
	}
	return UnknownGT
}
