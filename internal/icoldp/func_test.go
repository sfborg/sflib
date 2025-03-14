package icoldp_test

import (
	"strconv"
	"testing"

	"database/sql"

	"github.com/sfborg/sflib/pkg/coldp"
	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert := assert.New(t)
	testsStr := []struct {
		input         string
		valid, output bool
	}{
		{"smth", false, false},
		{"", false, false},
		{"true", true, true},
		{"TRUE", true, true},
		{"tRUe", true, true},
		{"T", true, true},
		{"yes", true, true},
		{"y", true, true},
		{"false", true, false},
		{"FAlse", true, false},
		{"no", true, false},
		{"N", true, false},
	}

	for _, v := range testsStr {
		res := coldp.ToBool(v.input)
		assert.Equal(v.output, res.Bool, v.input+"-string-bool")
		assert.Equal(v.valid, res.Valid, v.input+"-string-valid")
	}

	testsInt := []struct {
		input         int
		valid, output bool
	}{
		{1, true, true},
		{0, true, false},
		{2, false, false},
		{-1, false, false},
	}

	for _, v := range testsInt {
		res := coldp.ToBool(v.input)
		assert.Equal(v.output, res.Bool, strconv.Itoa(v.input)+"-int-bool")
		assert.Equal(v.valid, res.Valid, strconv.Itoa(v.input)+"-int-valid")
	}

	testsBool := []struct {
		input         bool
		valid, output bool
	}{
		{true, true, true},
		{false, true, false},
	}
	for _, v := range testsBool {
		res := coldp.ToBool(v.input)
		assert.Equal(v.output, res.Bool, strconv.FormatBool(v.input)+"-bool-bool")
		assert.Equal(v.valid, res.Valid, strconv.FormatBool(v.input)+"-bool-valid")
	}
}

func TestToInt(t *testing.T) {
	assert := assert.New(t)
	testsStr := []struct {
		input         string
		valid         bool
		output        int64
		outputDefault sql.NullInt64
	}{
		{"smth", false, 0, sql.NullInt64{}},
		{"", false, 0, sql.NullInt64{}},
		{"123", true, 123, sql.NullInt64{Int64: 123, Valid: true}},
		{"-123", true, -123, sql.NullInt64{Int64: -123, Valid: true}},
		{" 123 ", true, 123, sql.NullInt64{Int64: 123, Valid: true}},
	}

	for _, v := range testsStr {
		res := coldp.ToInt(v.input)
		assert.Equal(v.output, res.Int64, v.input+"-string-int")
		assert.Equal(v.valid, res.Valid, v.input+"-string-valid")
		assert.Equal(v.outputDefault, res, v.input+"-string-default")
	}

	testsInt := []struct {
		input         int
		valid         bool
		output        int64
		outputDefault sql.NullInt64
	}{
		{123, true, 123, sql.NullInt64{Int64: 123, Valid: true}},
		{-123, true, -123, sql.NullInt64{Int64: -123, Valid: true}},
		{0, true, 0, sql.NullInt64{Int64: 0, Valid: true}},
	}

	for _, v := range testsInt {
		res := coldp.ToInt(v.input)
		assert.Equal(v.output, res.Int64, strconv.Itoa(v.input)+"-int-int")
		assert.Equal(v.valid, res.Valid, strconv.Itoa(v.input)+"-int-valid")
		assert.Equal(v.outputDefault, res, strconv.Itoa(v.input)+"-int-default")
	}
}

func TestToFloat(t *testing.T) {
	assert := assert.New(t)
	testsStr := []struct {
		input         string
		valid         bool
		output        float64
		outputDefault sql.NullFloat64
	}{
		{"smth", false, 0, sql.NullFloat64{}},
		{"", false, 0, sql.NullFloat64{}},
		{"123.45", true, 123.45, sql.NullFloat64{Float64: 123.45, Valid: true}},
	}
	for _, v := range testsStr {
		res := coldp.ToFloat(v.input)
		assert.Equal(v.output, res.Float64, v.input)
		assert.Equal(v.valid, res.Valid, v.input)
		assert.Equal(v.outputDefault, res, v.input)
	}
	testsFloat := []struct {
		input         float64
		valid         bool
		output        float64
		outputDefault sql.NullFloat64
	}{
		{-123.45, true, -123.45, sql.NullFloat64{Float64: -123.45, Valid: true}},
	}
	for _, v := range testsFloat {
		res := coldp.ToFloat(v.input)
		assert.Equal(v.output, res.Float64, strconv.FormatFloat(v.input, 'f', -1, 64))
		assert.Equal(v.valid, res.Valid, strconv.FormatFloat(v.input, 'f', -1, 64))
		assert.Equal(v.outputDefault, res, strconv.FormatFloat(v.input, 'f', -1, 64))
	}
	testsInt := []struct {
		input         int
		valid         bool
		output        float64
		outputDefault sql.NullFloat64
	}{
		{123, true, 123, sql.NullFloat64{Float64: 123, Valid: true}},
	}
	for _, v := range testsInt {
		res := coldp.ToFloat(v.input)
		assert.Equal(v.output, res.Float64, strconv.Itoa(v.input))
		assert.Equal(v.valid, res.Valid, strconv.Itoa(v.input))
		assert.Equal(v.outputDefault, res, strconv.Itoa(v.input))
	}
}
