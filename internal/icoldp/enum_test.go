package icoldp_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestNomCodeNew(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, inp string
		out      coldp.NomCode
	}{
		{"bad", "smth", coldp.UnknownNomCode},
		{"zoo1", "zoological", coldp.Zoological},
		{"zoo2", "ICZN", coldp.Zoological},
		{"zoo3", "iczn", coldp.Zoological},
		{"zoo4", "Zoological", coldp.Zoological},
		{"bot1", "Botanical", coldp.Botanical},
		{"bot2", "ICN", coldp.Botanical},
		{"bot3", "ICNafp", coldp.Botanical},
		{"vir", "ICVCN", coldp.Virus},
		{"phyto", "ICPN", coldp.PhytoSociological},
		{"phyto", "ICNcp", coldp.Cultivars},
	}

	for _, v := range tests {
		res := coldp.NewNomCode(v.inp)
		assert.Equal(v.out, res, v.msg)
	}
}

func TestNomCodeString(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg string
		inp coldp.NomCode
		out string
	}{
		{"bad", coldp.UnknownNomCode, ""},
		{"vir", coldp.Virus, "VIRUS"},
		{"bot", coldp.NewNomCode("icn"), "BOTANICAL"},
	}

	for _, v := range tests {
		res := v.inp.ID()
		assert.Equal(v.out, res, v.msg)
	}
}

func TestNomRelType(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		inp coldp.NomRelType
		out string
	}{
		{coldp.SpellingCorrection, "SPELLING_CORRECTION"},
		{coldp.Basionym, "BASIONYM"},
		{coldp.BasedOn, "BASEDON"},
		{coldp.ReplacementName, "REPLACEMENT_NAME"},
		{coldp.ConservedNRT, "CONSERVED"},
		{coldp.LaterHomonym, "LATER_HOMONYM"},
		{coldp.Superfluous, "SUPERFLUOUS"},
		{coldp.Homotypic, "HOMOTYPIC"},
		{coldp.Type, "TYPE"},
	}

	for _, v := range testCases {
		res := v.inp.ID()
		assert.Equal(v.out, res)
	}
}

func TestNormRelTypeString(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		inp string
		out coldp.NomRelType
	}{
		{"spelling_correction", coldp.SpellingCorrection},
		{"BASIONYM", coldp.Basionym},
		{"  Basedon ", coldp.BasedOn},
		{"rePlacEMent_name", coldp.ReplacementName},
		{"Conserved", coldp.ConservedNRT},
		{"later homonym", coldp.LaterHomonym},
		{"superfluous  ", coldp.Superfluous},
		{"HOMOTYPIC_", coldp.Homotypic},
		{"TYPE", coldp.Type},
		{"Invalid", coldp.UnknownNomRelType},
		{"", coldp.UnknownNomRelType},
	}

	for _, v := range tests {
		res := coldp.NewNomRelType(v.inp)
		assert.Equal(v.out, res)
	}
}

func TestNamePartString(t *testing.T) {
	assert := assert.New(t)
	// Test namePartToString
	tests := []struct {
		inp coldp.NamePart
		out string
	}{
		{coldp.GenericNP, "GENERIC"},
		{coldp.InfragenericNP, "INFRAGENERIC"},
		{coldp.SpecificNP, "SPECIFIC"},
		{coldp.InfraspecificNP, "INFRASPECIFIC"},
	}

	for _, tc := range tests {
		res := tc.inp.ID()
		assert.Equal(tc.out, res)
	}
}

func TestNewNamePart(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		inp string
		out coldp.NamePart
	}{
		{"generic", coldp.GenericNP},
		{"Generic", coldp.GenericNP}, // Test case-insensitivity
		{"INFRAGENERIC", coldp.InfragenericNP},
		{"specific", coldp.SpecificNP}, // Test with extra spaces
		{"infraSpecific", coldp.InfraspecificNP},
		{"invalid", coldp.UnknownNP},
		{"", coldp.UnknownNP},
	}

	for _, tc := range tests {
		res := coldp.NewNamePart(tc.inp)
		assert.Equal(tc.out, res)
	}
}

func TestEnvironment(t *testing.T) {
	assert := assert.New(t)

	// Test envToString using Environment.String()
	tests := []struct {
		msg, inp string
		out      coldp.Environment
	}{
		{"brackish", "BRACKISH", coldp.Brackish},
		{"freshwater", "FRESHWATER", coldp.Freshwater},
		{"marine", "MARINE", coldp.Marine},
		{"terrestrial", "TERRESTRIAL", coldp.Terrestrial},
	}
	for _, v := range tests {
		res := v.out.ID()
		assert.Equal(v.inp, res, v.msg)
	}

	// Test stringToEnv using NewEnvironment()
	tests = []struct {
		msg, inp string
		out      coldp.Environment
	}{
		{"bad", "smth", coldp.UnknownEnv},
		{"brackish1", "brackish", coldp.Brackish},
		{"brackish2", "Brackish", coldp.Brackish},
		{"freshwater", "Freshwater", coldp.Freshwater},
		{"marine", "MARINE", coldp.Marine},
		{"terrestrial", "Terrestrial", coldp.Terrestrial},
	}

	for _, v := range tests {
		res := coldp.NewEnvironment(v.inp)
		assert.Equal(v.out, res, v.msg)
	}
}
