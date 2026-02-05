package idwca

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sfborg/sflib/pkg/dwca"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteMeta(t *testing.T) {
	dir := t.TempDir()
	a := &idwca{rootDir: dir}

	tests := []struct {
		name     string
		hasVern  bool
		hasDistr bool
		contains []string
		excludes []string
	}{
		{
			name:     "core only",
			hasVern:  false,
			hasDistr: false,
			contains: []string{
				`<archive`,
				`metadata="eml.xml"`,
				`rowType="http://rs.tdwg.org/dwc/terms/Taxon"`,
				`<location>Taxon.csv</location>`,
				`fieldsTerminatedBy=","`,
				`term="http://rs.tdwg.org/dwc/terms/taxonID"`,
				`term="http://rs.tdwg.org/dwc/terms/parentNameUsageID"`,
			},
			excludes: []string{
				"VernacularName",
				"Distribution",
			},
		},
		{
			name:     "core + vernacular + distribution",
			hasVern:  true,
			hasDistr: true,
			contains: []string{
				`<location>Taxon.csv</location>`,
				`<location>VernacularName.csv</location>`,
				`<location>Distribution.csv</location>`,
				`rowType="http://rs.gbif.org/terms/1.0/VernacularName"`,
				`rowType="http://rs.gbif.org/terms/1.0/Distribution"`,
				`<coreid index="0"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := dwca.BuildMeta(tt.hasVern, tt.hasDistr)
			err := a.WriteMeta(m)
			require.NoError(t, err)

			bs, err := os.ReadFile(filepath.Join(dir, "meta.xml"))
			require.NoError(t, err)

			xml := string(bs)
			for _, s := range tt.contains {
				assert.Contains(t, xml, s)
			}
			for _, s := range tt.excludes {
				assert.NotContains(t, xml, s)
			}
		})
	}
}
