package coldp

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/gnames/gnfmt/gncsv"
	csvConfig "github.com/gnames/gnfmt/gncsv/config"
	"github.com/gnames/gnlib"
	"github.com/sfborg/sflib/config"
)

// NormalizeHeaders attempts to normalize DarwinCore and ColDP terms to
// headers that correspond to ColDP.
func NormalizeHeaders(headers []string) map[string]int {
	res := map[string]int{}
	headers = gnlib.Map(headers, func(s string) string {
		s = strings.ToLower(s)
		els := strings.Split(s, ":")
		return els[len(els)-1]
	})

	for i, v := range headers {
		switch v {
		case "id", "taxonid":
			res["id"] = i
			res["nameid"] = i
		case "parentid", "parentnameusageid":
			res["parentid"] = i
		case "scientificname":
			res["scientificname"] = i
		case "authorship", "scientificnameauthorship":
			res["authorship"] = i
		case "notho":
			res["notho"] = i
		case "rank", "taxonrank":
			res["rank"] = i
		case "genus", "genericname":
			res["genericname"] = i
		case "infragenericepithet", "subgenus", "infragenus":
			res["infragenericepithet"] = i
		case "specificepithet", "species":
			res["specificepithet"] = i
		case "infraspecificepithet", "infraspecies":
			res["infraspecificepithet"] = i
		case "code", "nomenclaturalcode", "nomcode":
			res["code"] = i
		case "remarks", "taxonremarks":
			res["remarks"] = i
		case "link", "url":
			res["link"] = i
		default:
			res[v] = i
		}
	}
	return res
}

func ReadJSON[T DataLoader](
	path string,
	chIn chan T) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	slog.Info("Parsing references JSON-CSL file")
	var cslRec CSLRecords
	err = json.Unmarshal(byteValue, &cslRec)
	if err != nil {
		return err
	}
	for _, csl := range cslRec {
		chIn <- csl.ToReference().(T)
	}
	close(chIn)
	return nil
}

func ReadJSONL[T DataLoader](
	path string,
	chIn chan T) error {
	f, err := os.Open(path) // Replace "my_file.txt" with your file name
	if err != nil {
		return err
	}
	defer f.Close() // Ensure the file is closed when done

	slog.Info("Reading references from JSON-CSL file")

	// Create a new scanner
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file
	for scanner.Scan() {
		bs := scanner.Bytes()
		var csl CSL
		err = json.Unmarshal(bs, &csl)
		if err != nil {
			return err
		}
		chIn <- csl.ToReference().(T)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	close(chIn)
	return nil
}

func Read[T DataLoader](
	cfg config.ConfigCoLDP,
	path string,
	chOut chan T) error {
	chIn := make(chan []string)
	var wg sync.WaitGroup
	wg.Add(1)

	opts := []csvConfig.Option{
		csvConfig.OptPath(path),
		csvConfig.OptBadRowMode(cfg.BadRow),
	}
	csvCfg, err := csvConfig.New(opts...)
	headers := gnlib.Map(csvCfg.Headers, func(s string) string {
		s = strings.ToLower(s)
		els := strings.Split(s, ":")
		return els[len(els)-1]
	})

	var colSep string
	switch csvCfg.ColSep {
	case '\t':
		colSep = "'\\t'"
	default:
		colSep = "'" + string(csvCfg.ColSep) + "'"
	}
	slog.Info("Processing CSV", "col-sep", colSep, "file", filepath.Base(path))

	go func() {
		defer wg.Done()
		for row := range chIn {
			var t T
			dl, warn := t.Load(headers, row)
			t = dl.(T)
			if warn != nil {
				slog.Warn("Cannot read row", "warn", warn, "row", row)
				continue
			}
			chOut <- t
		}
	}()

	if err != nil {
		return err
	}

	csv := gncsv.New(csvCfg)
	_, err = csv.Read(context.Background(), chIn)
	if err != nil {
		return err
	}
	close(chIn)
	wg.Wait()

	return nil
}

type FieldNumberWarning struct {
	FieldsNum, RowFieldsNum int
	Row                     []string
	Message                 string
}

func (w *FieldNumberWarning) Error() string {
	return w.Message
}

// ToInt attempts to convert a value to a sql.NullInt64,
// handling both string and int input types.
func ToInt[T string | int](val T) sql.NullInt64 {
	var res sql.NullInt64

	switch v := any(val).(type) {
	case string:
		s := strings.ToLower(v)
		s = strings.TrimSpace(s)
		i, err := strconv.Atoi(s)
		if err != nil {
			return res
		}
		res.Int64 = int64(i)
		res.Valid = true
	case int:
		res.Int64 = int64(v)
		res.Valid = true
	default:
		res.Valid = false
	}

	return res
}

// ToBool attempts to convert a value to a sql.NullBool,
// handling both string and bool input types.
func ToBool[T string | bool | int](val T) sql.NullBool {
	var res sql.NullBool
	res.Valid = true

	switch v := any(val).(type) {
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			res.Valid = false
			return res
		}

		s = strings.ToLower(s)
		b := s[0]
		switch b {
		case '1', 'y', 't':
			res.Bool = true
		case '0', 'n', 'f':
			res.Bool = false
		default:
			res.Bool = false
			res.Valid = false
		}
	case int:
		if v == 1 {
			res.Bool = true
		} else if v == 0 {
			res.Bool = false
		} else {
			res.Valid = false
		}
	case bool:
		res.Bool = v
	default:
		res.Valid = false
	}

	return res
}

// ToFloat attempts to convert a value to a sql.NullFloat64,
// handling string, float64, and int input types.
func ToFloat[T string | float64 | int](val T) sql.NullFloat64 {
	var res sql.NullFloat64

	switch v := any(val).(type) {
	case string:
		s := strings.ToLower(v)
		s = strings.TrimSpace(s)
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return res
		}
		res.Float64 = f
		res.Valid = true
	case float64:
		res.Float64 = v
		res.Valid = true
	case int:
		res.Float64 = float64(v)
		res.Valid = true
	default:
		res.Valid = false
	}

	return res
}

func RowToMap(headers, row []string) (map[string]string, error) {
	diff := len(headers) - len(row)
	var warning error
	if diff > 0 {
		for range diff {
			row = append(row, "")
		}
		warning = &FieldNumberWarning{
			FieldsNum:    len(headers),
			RowFieldsNum: len(row),
			Row:          row,
			Message:      fmt.Sprintf("not enough fields, filled with empty strings"),
		}
	}
	if diff < 0 {
		warning = &FieldNumberWarning{
			FieldsNum:    len(headers),
			RowFieldsNum: len(row),
			Row:          row,
			Message:      fmt.Sprintf("too many fields, extras will be ignored"),
		}
	}

	res := make(map[string]string)
	for i := range headers {
		if i < len(row) { // Prevent index out of range
			fld := strings.ToLower(headers[i])
			fld = strings.ReplaceAll(fld, "_", "")
			res[headers[i]] = row[i]
		}
	}

	return res, warning
}

// ToStr normalizes enumerated string IDs to 'normal' strings.
// For example 'PROVISIONALLY_ACCEPTED' becomes
// 'provisionally accepted'.
func ToStr(s string) string {
	s = strings.ToLower(s)
	return strings.ReplaceAll(s, "_", " ")
}

func ToTitleCase(s string) string {
	words := strings.Split(s, "_")
	var result strings.Builder
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			result.WriteString(string(runes))
			if i < len(words)-1 {
				result.WriteString(" ")
			}
		}
	}
	return result.String()
}

func GetEnvironments(env string) []Environment {
	envs := strings.Split(env, ",")
	envs = gnlib.Map(envs, func(s string) string {
		return strings.TrimSpace(s)
	})
	var res []Environment
	for _, v := range envs {
		env := NewEnvironment(v)
		if env == UnknownEnv {
			slog.Warn("Unknown environment", "env", v)
			continue
		}
		res = append(res, env)
	}
	return res
}
