module github.com/sfborg/sflib

go 1.25.1

// Local fork with a fix for SQLite shorthand REFERENCES introspection.
// Swap for the pushed fork reference once available upstream or on GitHub.
replace ariga.io/atlas => github.com/dimus/atlas v0.0.0-20260512202807-c70f08441137

require (
	ariga.io/atlas v1.2.1
	github.com/dustin/go-humanize v1.0.1
	github.com/gnames/gnfmt v0.6.5
	github.com/gnames/gnlib v0.64.0
	github.com/gnames/gnparser v1.15.0
	github.com/gnames/gnsys v0.4.4
	github.com/gnames/gnuuid v0.2.0
	github.com/goccy/go-yaml v1.19.2
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/sync v0.20.0
	modernc.org/sqlite v1.51.0
)

require (
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/cheggaaa/pb/v3 v3.1.7 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/gnames/organizer v0.1.1 // indirect
	github.com/gnames/tribool v0.1.1 // indirect
	github.com/go-openapi/inflect v0.21.5 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/zclconf/go-cty v1.18.1 // indirect
	github.com/zclconf/go-cty-yaml v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20260529124908-c761662dc8c9 // indirect
	golang.org/x/mod v0.36.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/tools v0.45.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.72.5 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)
