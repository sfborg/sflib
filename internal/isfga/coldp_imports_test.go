package isfga_test

import (
	"os"
	"testing"

	"github.com/gnames/coldp/ent/coldp"
	"github.com/sfborg/sflib/internal/isfga"
	"github.com/sfborg/sflib/pkg/arch"
	"github.com/stretchr/testify/assert"
)

var (
	a       arch.SFGA
	testDir string
)

func TestMain(m *testing.M) {
	setupGlobal()
	code := m.Run() // Run all tests
	teardownGlobal()
	os.Exit(code)
}

func setupGlobal() {
	var err error
	a = isfga.New()
	testDir, err = os.MkdirTemp("", "sfga-test")
	if err != nil {
		panic(err)
	}
	err = a.Create(testDir)
	if err != nil {
		panic(err)
	}
	_, err = a.Connect()
}

func teardownGlobal() {
	var err error
	err = a.Close()
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(testDir)
	if err != nil {
		panic(err)
	}
}

func TestInsertAuthors(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	authors := []coldp.Author{
		{
			ID:                 "au1",
			SourceID:           "src1",
			AlternativeID:      "alt1",
			Given:              "John",
			Family:             "Doe",
			Suffix:             "Jr",
			AbbreviationBotany: "J.Doe",
			AlternativeNames:   "Johnny Doe",
			Sex:                coldp.NewSex("male"),
			Country:            "USA",
		},
	}
	err := a.InsertAuthors(authors)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from author").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertDistributions(t *testing.T) {
	assert := assert.New(t)
	distr := []coldp.Distribution{
		{
			TaxonID:     "tx1",
			SourceID:    "src1",
			Area:        "area1",
			AreaID:      "ar1",
			Gazetteer:   coldp.NewGazetteerEnt("text"),
			Status:      coldp.NewDistrStatus("invasive"),
			ReferenceID: "ref1",
			Remarks:     "rem1",
			Modified:    "mod1",
			ModifiedBy:  "modby1",
		},
		{
			TaxonID:     "tx2",
			SourceID:    "src2",
			Area:        "area2",
			AreaID:      "ar2",
			Gazetteer:   coldp.NewGazetteerEnt("text"),
			Status:      coldp.NewDistrStatus("invasive"),
			ReferenceID: "ref2",
			Remarks:     "rem2",
			Modified:    "mod2",
			ModifiedBy:  "modby2",
		},
	}
	err := a.InsertDistributions(distr)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from distribution").Scan(&res)
	assert.Nil(err)
	assert.Equal(2, res)
}

func TestInsertMedia(t *testing.T) {
	assert := assert.New(t)
	media := []coldp.Media{
		{
			TaxonID:  "tx1",
			SourceID: "src1",
			URL:      "url1",
			Type:     "type1",
			Format:   "format1",
			Title:    "title1",
			Created:  "created1",
			Creator:  "creator1",
			License:  "license1",
			Link:     "link1",
			Remarks:  "remarks1",
		},
		{
			TaxonID:  "tx2",
			SourceID: "src2",
			URL:      "url2",
			Type:     "type2",
			Format:   "format2",
			Title:    "title2",
			Created:  "created2",
			Creator:  "creator2",
		},
	}
	err := a.InsertMedia(media)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from media").Scan(&res)
	assert.Nil(err)
	assert.Equal(2, res)
}

func TestInsertMeta(t *testing.T) {
	assert := assert.New(t)
	meta := coldp.Meta{
		DOI:         "doi1",
		Title:       "title1",
		Alias:       "alias1",
		Description: "desc1",
		Issued:      "issued1",
		Version:     "version1",
		Keywords:    []string{"key1", "key2"},
		License:     "license1",
		URL:         "url1",
		Logo:        "logo1",
		Label:       "label1",
		Citation:    "citation1",
		Private:     true,
		Contact:     &coldp.Actor{Given: "John", Family: "Doe"},
		Publisher:   &coldp.Actor{Given: "John", Family: "Doe"},
		Editors: []coldp.Actor{
			{Given: "John", Family: "Doe"},
			{Given: "Peter", Family: "Pan"},
			{Given: "Ben", Family: "Doe"},
		},
		Creators: []coldp.Actor{
			{Given: "John", Family: "Doe"},
			{Given: "Peter", Family: "Pan"},
		},
		Contributors: []coldp.Actor{
			{Given: "Jane", Family: "Doe", City: "Urbana", State: "IL"},
			{Given: "Peter", Family: "Pan"},
		},
	}
	err := a.InsertMeta(&meta)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from metadata").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
	err = db.QueryRow("select count(*) from publisher").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
	err = db.QueryRow("select count(*) from contact").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
	err = db.QueryRow("select count(*) from editor").Scan(&res)
	assert.Nil(err)
	assert.Equal(3, res)
	err = db.QueryRow("select count(*) from creator").Scan(&res)
	assert.Nil(err)
	assert.Equal(2, res)
	err = db.QueryRow("select count(*) from contributor").Scan(&res)
	assert.Nil(err)
	assert.Equal(2, res)
}

func TestInsertNameRelations(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	rel := []coldp.NameRelation{
		{
			NameID:        "nm1",
			SourceID:      "src1",
			Type:          coldp.NewNomRelType("homonym"),
			RelatedNameID: "rel1",
			ReferenceID:   "ref1",
		},
		{
			NameID:        "nm2",
			SourceID:      "src2",
			Type:          coldp.NewNomRelType("basionym"),
			RelatedNameID: "rel2",
		},
	}
	db := a.Db()
	_, err := db.Exec("delete from name_relation")
	assert.Nil(err)

	err = a.InsertNameRelations(rel)
	assert.Nil(err)
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from name_relation").Scan(&res)
	assert.Nil(err)
	assert.Equal(2, res)
}

func TestInsertNameUsage(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	nu := []coldp.NameUsage{
		{
			ID:                   "123",
			AlternativeID:        "abc",
			NameAlternativeID:    "234",
			LocalID:              "local1",
			GlobalID:             "global1",
			SourceID:             "source1",
			ParentID:             "321",
			BasionymID:           "33",
			TaxonomicStatus:      coldp.NewTaxonomicStatus("accepted"),
			ScientificName:       "Bubo bubo",
			Authorship:           "L.",
			ScientificNameString: "Bubo bubo L.",
			Rank:                 coldp.NewRank("species"),
			GenericName:          "Bubo",
			SpecificEpithet:      "bubo",
			BasionymAuthorship:   "L.",
			NameReferenceID:      "123",
			PublishedInYear:      "1758",
			PublishedInPage:      "123",
			PublishedInPageLink:  "http://example.org",
			Code:                 coldp.NewNomCode("zoological"),
			NameStatus:           coldp.NewNomStatus("acceptable"),
			ReferenceID:          "123",
			Scrutinizer:          "James Bond",
			ScrutinizerDate:      "2020-12-12",
			Extinct:              coldp.ToBool("false"),
			Environment:          coldp.NewEnvironment("terrestrial"),
			Link:                 "http://example.org",
			NameRemarks:          "name rem",
			Remarks:              "rem",
			Modified:             "2023-01-01",
			ModifiedBy:           "Lin.",
		},
	}
	// clean up taxon and name
	db := a.Db()
	_, err := db.Exec("delete from name")
	_, err = db.Exec("delete from taxon")
	assert.Nil(err)

	err = a.InsertNameUsages(nu)
	assert.Nil(err)
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from name").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
	err = db.QueryRow("select count(*) from taxon").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertNames(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	name := []coldp.Name{
		{
			ID:                   "123",
			AlternativeID:        "321",
			SourceID:             "aa",
			BasionymID:           "555",
			ScientificName:       "Bubo bubo",
			Authorship:           "L.",
			ScientificNameString: "Bubo bubo L.",
			Rank:                 coldp.NewRank("species"),
			Genus:                "Bubo",
			SpecificEpithet:      "bubo",
			Code:                 coldp.NewNomCode("zoological"),
			Status:               coldp.NewNomStatus("acceptable"),
			ReferenceID:          "123",
			PublishedInYear:      "1758",
			PublishedInPage:      "123",
			PublishedInPageLink:  "http://example.org",
			Link:                 "http://example.org",
			Remarks:              "rem",
			Modified:             "2022-02-01",
			ModifiedBy:           "John",
		},
	}
	// truncate name table other tests data
	db := a.Db()
	_, err := db.Exec("delete from name")
	assert.Nil(err)

	err = a.InsertNames(name)
	assert.Nil(err)
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from name").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertReferences(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	ref := []coldp.Reference{
		{
			ID:        "ref1",
			SourceID:  "src1",
			Citation:  "citation1",
			Type:      coldp.NewReferenceType("book"),
			Author:    "author1",
			Title:     "title1",
			Issued:    "issued1",
			Publisher: "publisher1",
		},
	}
	err := a.InsertReferences(ref)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from reference").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertSpeciesEstimate(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	se := []coldp.SpeciesEstimate{
		{
			TaxonID:     "tx1",
			SourceID:    "src1",
			Estimate:    coldp.ToInt("12022928"),
			Type:        coldp.NewEstimateType("count"),
			ReferenceID: "ref1",
			Remarks:     "rem1",
			Modified:    "mod1",
			ModifiedBy:  "modby1",
		},
	}
	err := a.InsertSpeciesEstimates(se)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from species_estimate").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertSpeciesInteractions(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	si := []coldp.SpeciesInteraction{
		{
			TaxonID:                    "tx1",
			RelatedTaxonID:             "tx2",
			SourceID:                   "src1",
			RelatedTaxonScientificName: "Bubo bubo",
			Type:                       coldp.NewSpInteractionType("parasite"),
			ReferenceID:                "ref1",
			Remarks:                    "rem1",
			Modified:                   "mod1",
			ModifiedBy:                 "modby1",
		},
	}
	err := a.InsertSpeciesInteractions(si)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from species_interaction").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertSynonyms(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	syn := []coldp.Synonym{
		{
			ID:            "syn1",
			TaxonID:       "tx1",
			SourceID:      "src1",
			NameID:        "nm1",
			NamePhrase:    "phrase1",
			AccordingToID: "acc1",
			Status:        coldp.NewTaxonomicStatus("synonym"),
			ReferenceID:   "ref1",
			Link:          "link1",
			Remarks:       "rem1",
		},
	}
	err := a.InsertSynonyms(syn)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from synonym").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertTxConcRelations(t *testing.T) {
	assert := assert.New(t)
	assert.True(a.Ping())
	tcr := []coldp.TaxonConceptRelation{
		{
			TaxonID:        "tx1",
			RelatedTaxonID: "tx2",
			SourceID:       "src1",
			Type:           coldp.NewTaxonConceptRelType("equals"),
			ReferenceID:    "ref1",
			Remarks:        "rem1",
			Modified:       "mod1",
		},
	}
	err := a.InsertTaxonConceptRelations(tcr)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from taxon_concept_relation").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertTxProperties(t *testing.T) {
	assert := assert.New(t)
	tp := []coldp.TaxonProperty{
		{
			TaxonID:     "tx1",
			SourceID:    "src1",
			Property:    "prop",
			ReferenceID: "ref1",
		},
	}
	err := a.InsertTaxonProperties(tp)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from taxon_property").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertTaxon(t *testing.T) {
	assert := assert.New(t)
	tx := []coldp.Taxon{
		{
			ID:                  "123",                               // string
			AlternativeID:       "321",                               // string
			LocalID:             "111",                               // string
			GlobalID:            "222",                               // string
			SourceID:            "123",                               // string
			ParentID:            "543",                               // string
			NameID:              "333",                               // string
			AccordingToID:       "555",                               // string
			AccordingToPage:     "66",                                // string
			AccordingToPageLink: "http://example.org",                // string
			Scrutinizer:         "Li Xi",                             // string
			ScrutinizerID:       "88",                                // string
			ScrutinizerDate:     "2020-02-21",                        // string
			Provisional:         coldp.ToBool("false"),               // sql.NullBool
			ReferenceID:         "44",                                // string
			Extinct:             coldp.ToBool("false"),               // sql.NullBool
			Environment:         coldp.NewEnvironment("terrestrial"), // Environment
			Species:             "Bubo",                              // string
			Family:              "Strigidae",                         // string
			Order:               "Strigiformes",                      // string
			Class:               "Aves",                              // string
			Kingdom:             "Animalia",                          // string
			Link:                "http://example.org",                // string
			Remarks:             "rem",                               // string
			Modified:            "2022-02-03",                        // string
			ModifiedBy:          "Xin Min",                           // string
		},
	}

	// we need to trim taxon table, because InsertNameUsage also inserts into
	// it
	db := a.Db()
	_, err := db.Exec("delete from taxon")
	assert.Nil(err)

	err = a.InsertTaxa(tx)
	assert.Nil(err)
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from taxon").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertTreatment(t *testing.T) {
	assert := assert.New(t)
	tr := []coldp.Treatment{
		{
			TaxonID:  "tr1",
			SourceID: "src1",
			Document: "ref1",
			Format:   "format1",
		},
	}
	err := a.InsertTreatments(tr)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from treatment").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertTypeMaterials(t *testing.T) {
	assert := assert.New(t)
	tm := []coldp.TypeMaterial{
		{
			ID:          "tx1",
			SourceID:    "src1",
			NameID:      "nm1",
			Citation:    "cit1",
			Status:      coldp.NewTypeStatus("holotype"),
			ReferenceID: "ref1",
			Locality:    "loc1",
			Latitude:    coldp.ToFloat("33.3"),
			Longitude:   coldp.ToFloat("33.3"),
			Altitude:    coldp.ToInt("44"),
		},
	}
	err := a.InsertTypeMaterials(tm)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from type_material").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}

func TestInsertVernaculars(t *testing.T) {
	assert := assert.New(t)
	vrn := []coldp.Vernacular{
		{
			SourceID:    "src1",
			TaxonID:     "tx1",
			Name:        "name1",
			Language:    "lang1",
			Country:     "country1",
			ReferenceID: "ref1",
		},
	}
	err := a.InsertVernaculars(vrn)
	assert.Nil(err)
	db := a.Db()
	assert.NotNil(db)

	var res int
	err = db.QueryRow("select count(*) from vernacular").Scan(&res)
	assert.Nil(err)
	assert.Equal(1, res)
}
