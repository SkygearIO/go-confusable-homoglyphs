package confusablehomoglyphs

import (
	"fmt"
	"testing"
)

func TestIsMixedScript(t *testing.T) {
	cases := []struct {
		str            string
		allowedAliases []string
		isMixedScript  bool
	}{
		{"Abç", nil, false},
		{"ρτ.τ", nil, false},
		{"ρτ.τ", []string{}, true},
		{"Alloτ", nil, true},
	}

	for _, c := range cases {
		fmt.Println("---------------------- case : ", c.str)
		isMixedScript := IsMixedScript(c.str, c.allowedAliases)
		if isMixedScript != c.isMixedScript {
			t.Errorf("unexpected isMixedScript, string: %v, expected: %v, actual: %v\n", c.str, c.isMixedScript, isMixedScript)
		}
	}
}

func TestIsConfusable(t *testing.T) {
	cases := []struct {
		str              string
		preferredAliases []string
		isConfusable     bool
		checkingFunction func([]ConfusableResult)
	}{
		{"paρa", []string{"latin"}, true, func(r []ConfusableResult) {
			if r[0].Character != 'ρ' {
				t.Errorf("unexpected confusable result: %v\n", r)
			}
		}},
		{"paρa", []string{"greek"}, true, func(r []ConfusableResult) {
			if r[0].Character != 'p' {
				t.Errorf("unexpected confusable result: %v\n", r)
			}
		}},
		{"Abç", []string{"latin"}, false, func(r []ConfusableResult) {}},
		{"AlloΓ", []string{"latin"}, false, func(r []ConfusableResult) {}},
		{"ρττ", []string{"greek"}, false, func(r []ConfusableResult) {}},
		{"ρτ.τ", []string{"greek", "common"}, false, func(r []ConfusableResult) {}},
		{"ρττp", nil, true, func(r []ConfusableResult) {
			if r[0].Character != 'ρ' ||
				r[0].Alias != "GREEK" ||
				r[0].Homoglyphs[0].C != "p" ||
				r[0].Homoglyphs[0].N != "LATIN SMALL LETTER P" {
				t.Errorf("unexpected confusable result: %v\n", r)
			}
		}},
	}

	for _, c := range cases {
		confusableResult := IsConfusable(c.str, false, c.preferredAliases)
		isConfusable := len(confusableResult) > 0
		if isConfusable != c.isConfusable {
			t.Errorf("unexpected isConfusable, string: %v, expected: %v, actual: %v\n", c.str, c.isConfusable, isConfusable)
		}
		c.checkingFunction(confusableResult)
	}
}

func TestSetConfusableToLatin(t *testing.T) {
	cases := []struct {
		str              string
		preferredAliases []string
		isConfusable     bool
		checkingFunction func(string)
	}{
		{"paρa", []string{"latin"}, true, func(r string) {
			if r[2] != 'p' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"ρhishlabs", []string{"latin"}, true, func(r string) {
			if r[0] != 'p' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"Abç", []string{"latin"}, true, func(r string) {
			if r[2] == 'ç' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"Alloα", []string{"latin"}, true, func(r string) {
			if r[4] != 'a' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"ᑲankofamericα", []string{"latin"}, true, func(r string) {
			if r[0] != 'b' && r[12] != 'a' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"f1ipcaгt", []string{"latin"}, true, func(r string) {
			if r[1] != 'l' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"νegetables", []string{"latin"}, true, func(r string) {
			if r[0] != 'v' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"ոews", []string{"latin"}, true, func(r string) {
			if r[0] != 'n' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"amaz໐n", []string{"latin"}, true, func(r string) {
			if r[4] != 'o' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"AlloΓ", []string{"latin"}, true, func(r string) {
			if r != "Allo" {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
		{"ρτ.τ", []string{"greek", "common"}, true, func(r string) {}},
		{"ρττ", []string{"greek"}, true, func(r string) {}},
		{"Alloτ", []string{"latin"}, true, func(r string) {}},
		{"ρττp", []string{"latin"}, true, func(r string) {
			if r[0] != 'p' && r[1] != 'T' && r[2] != 't' && r[3] != 'T' {
				t.Errorf("unexpected latin string: %v\n", r)
			}
		}},
	}

	for _, c := range cases {
		latinResult := SetConfusableToLatin(c.str, c.preferredAliases)
		isConfusable := len(latinResult) > 0
		fmt.Println("ConfusableToLatin String =>  ", latinResult)

		if isConfusable != c.isConfusable {
			t.Errorf("unexpected isConfusable, string: %v, expected: %v, actual: %v\n", c.str, c.isConfusable, isConfusable)
		}
		c.checkingFunction(latinResult)
	}
}

func TestDangerous(t *testing.T) {
	cases := []struct {
		str              string
		preferredAliases []string
		isDangerous      bool
	}{
		{"Allo", []string{}, false},
		{"AlloΓ", []string{"latin"}, false},
		{"Alloρ", []string{}, true},
		{"AlaskaJazz", []string{}, false},
		{"ΑlaskaJazz", []string{}, true},
	}

	for _, c := range cases {
		isDangerous := IsDangerous(c.str, c.preferredAliases)
		if isDangerous != c.isDangerous {
			t.Errorf("unexpected isDangerous, string: %v, expected: %v, actual: %v\n", c.str, c.isDangerous, isDangerous)
		}
	}
}
