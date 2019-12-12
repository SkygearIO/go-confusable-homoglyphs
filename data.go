package confusablehomoglyphs

import (
	"encoding/json"
)

var confusablesData = map[string][]Homoglyph{}
var categoryData = CategoryData{}

type CategoryData struct {
	ISO15924aliases  []string `json:"iso_15924_aliases"`
	Categories       []string `json:"categories"`
	CodePointsRanges [][]int  `json:"code_points_ranges"`
}

type Homoglyph struct {
	C string `json:"c"`
	N string `json:"n"`
}

func init() {
	var err error
	err = json.Unmarshal([]byte(confusablesJSONTXT), &confusablesData)
	if err != nil {
		panic("failed to parse confusables json")
	}

	err = json.Unmarshal([]byte(categoriesJSONTXT), &categoryData)
	if err != nil {
		panic("failed to parse categories json")
	}
}
