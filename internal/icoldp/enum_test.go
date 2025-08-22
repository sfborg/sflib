package icoldp_test

import (
	"testing"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

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

// TestNewNomStatus tests the NewNomStatus function
func TestNewNomStatus(t *testing.T) {
	assert := assert.New(t)

	// Disable slog warnings during test
	tests := []struct {
		input    string
		expected coldp.NomStatus
	}{
		// Empty string
		{"", coldp.Established},

		// Established cases
		{"nomen validum", coldp.Established},
		{"NOMEN_VALIDUM", coldp.Established},
		{"available", coldp.Established},
		{"ESTABLISHED", coldp.Established},
		{"Established", coldp.Established},

		// NotEstablished cases
		{"nom. inval.", coldp.NotEstablished},
		{"nomen invalidum", coldp.NotEstablished},
		{"unavailable", coldp.NotEstablished},
		{"not established", coldp.NotEstablished},
		{"NOT_ESTABLISHED", coldp.NotEstablished},

		// Acceptable cases
		{"nomen legitimum", coldp.Acceptable},
		{"potentially valid", coldp.Acceptable},
		{"ACCEPTABLE", coldp.Acceptable},
		{"acceptable", coldp.Acceptable},

		// Unacceptable cases
		{"nom. illeg.", coldp.Unacceptable},
		{"nomen illegitimum", coldp.Unacceptable},
		{"NOME_ILEGITIMO", coldp.Unacceptable},
		{"objectively invalid", coldp.Unacceptable},
		{"unacceptable", coldp.Unacceptable},
		{"nudum", coldp.Unacceptable},

		// Conserved cases
		{"nom. cons.", coldp.Conserved},
		{"nomen conservandum", coldp.Conserved},
		{"conserved name", coldp.Conserved},
		{"CONSERVED", coldp.Conserved},
		{"NOME_CORRETO_VIA_CONSERVACAO", coldp.Conserved},

		// Rejected cases
		{"nom. rej.", coldp.Rejected},
		{"NOME_REJEITADO", coldp.Rejected},
		{"nomen rejiciendum", coldp.Rejected},
		{"rejected", coldp.Rejected},

		// Doubtful cases
		{"nom. dub.", coldp.Doubtful},
		{"nomen dubium", coldp.Doubtful},
		{"doubtful", coldp.Doubtful},
		{"dubium", coldp.Doubtful},

		// Manuscript cases
		{"manuscript name", coldp.Manuscript},
		{"manuscript", coldp.Manuscript},
		{"provisorium", coldp.Manuscript},

		// Chresonym cases
		{"chresonym", coldp.Chresonym},
		{"CHRESONYM", coldp.Chresonym},

		// URL cases
		{"http://example.com/established", coldp.Established},
		{"https://vocab.org/nomen_validum", coldp.Established},
		{"http://terms.gbif.org/doubtful", coldp.Doubtful},

		// Alternativum (maps to UnknownNomStatus)
		{"alternativum", coldp.UnknownNomStatus},

		// Unknown cases
		{"NOME_APLICACAO_INCERTA", coldp.UnknownNomStatus},
		{"NOME_NAO_EFETIVAMENTE_PUBLICADO", coldp.UnknownNomStatus},
		{"NOME_NAO_VALIDAMENTE_PUBLICADO", coldp.UnknownNomStatus},
		{"VARIANTE_ORTOGRAFICA", coldp.UnknownNomStatus},
		{"NOME_LEGITIMO_MAS_INCORRETO", coldp.UnknownNomStatus},
		{"NOME_MAL_APLICADO", coldp.UnknownNomStatus},
		{"VARIANTE_ORTOGRAFICA", coldp.UnknownNomStatus},
		{"NOME_CORRETO", coldp.UnknownNomStatus},
		{"unknown_status", coldp.UnknownNomStatus},
		{"invalid_input", coldp.UnknownNomStatus},
	}

	for _, test := range tests {
		result := coldp.NewNomStatus(test.input)
		assert.Equal(test.expected, result, test.input)
	}
}
