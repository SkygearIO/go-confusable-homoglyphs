package confusablehomoglyphs

import (
	"reflect"
	"sort"
	"testing"
)

func TestAliasesCategories(t *testing.T) {
	cases := []struct {
		char     rune
		alias    string
		category string
	}{
		{'A', "LATIN", "L"},
		{'τ', "GREEK", "L"},
		{'-', "COMMON", "Pd"},
	}

	for _, c := range cases {
		alias, cat := AliasesCategories(c.char)
		if alias != c.alias {
			t.Errorf("unexpected alias, expected: %v, actual: %v\n", c.alias, alias)
		}

		if cat != c.category {
			t.Errorf("unexpected category, expected: %v, actual: %v\n", c.category, cat)
		}
	}
}

func TestAlias(t *testing.T) {
	cases := []struct {
		char  rune
		alias string
	}{
		{'A', "LATIN"},
		{'τ', "GREEK"},
		{'-', "COMMON"},
	}

	for _, c := range cases {
		alias := Alias(c.char)
		if alias != c.alias {
			t.Errorf("unexpected alias, expected: %v, actual: %v\n", c.alias, alias)
		}
	}
}

func TestCategory(t *testing.T) {
	cases := []struct {
		char     rune
		category string
	}{
		{'A', "L"},
		{'τ', "L"},
		{'-', "Pd"},
	}

	for _, c := range cases {
		category := Category(c.char)
		if category != c.category {
			t.Errorf("unexpected category, expected: %v, actual: %v\n", c.category, category)
		}
	}
}

func TestUniqueAliases(t *testing.T) {
	cases := []struct {
		str     string
		aliases []string
	}{
		{"ABC", []string{"LATIN"}},
		{"ρAτ-", []string{"GREEK", "LATIN", "COMMON"}},
	}

	for _, c := range cases {
		aliases := UniqueAliases(c.str)
		if !sameArray(aliases, c.aliases) {
			t.Errorf("unexpected unique aliases, expected: %v, actual: %v\n", c.aliases, aliases)
		}
	}
}

func sameArray(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	ss1 := sort.StringSlice(arr1)
	ss1.Sort()
	ss2 := sort.StringSlice(arr2)
	ss2.Sort()

	return reflect.DeepEqual(ss1, ss2)
}
