// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"text/template"
)

func main() {
	confusables, err := ioutil.ReadFile("./tools/confusables.json")
	if err != nil {
		log.Fatal(err)
	}

	categories, err := ioutil.ReadFile("./tools/categories.json")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("data_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	packageTemplate.Execute(f, struct {
		ConfusablesJSON string
		CategoriesJSON  string
	}{
		ConfusablesJSON: strconv.Quote(string(confusables)),
		CategoriesJSON:  strconv.Quote(string(categories)),
	})
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package confusablehomoglyphs

const confusablesJSONTXT = {{ .ConfusablesJSON }}

const categoriesJSONTXT = {{ .CategoriesJSON }}

`))
