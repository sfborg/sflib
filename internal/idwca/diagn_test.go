package idwca_test

// func TestHierarchy(t *testing.T) {
// 	assert := assert.New(t)
// 	var err error
// 	dir := filepath.Join(testDir, "hier")
// 	err = gnsys.MakeDir(dir)
// 	assert.Nil(err)
//
// 	tests := []struct {
// 		msg, file string
// 		hierType  diagn.HierType
// 	}{
// 		{"flat", "flat.tar.gz", diagn.HierFlat},
// 	}
//
// 	for _, v := range tests {
// 		err := gnsys.CleanDir(dir)
// 		assert.Nil(err)
//
// 		path := filepath.Join("../../testdata/dwca/diagn/hierarchy/", v.file)
// 		a := idwca.New()
// 		err = a.Import(path, dir)
// 		assert.Nil(err)
// 	}
//
// }
