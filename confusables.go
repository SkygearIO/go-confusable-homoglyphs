package confusablehomoglyphs

import (
	"strings"
)

type ConfusableResult struct {
	Character  rune        `json:"character"`
	Alias      string      `json:"alias"`
	Homoglyphs []Homoglyph `json:"homoglyphs"`
}

func IsMixedScript(str string, allowedAliases []string) bool {
	if allowedAliases == nil {
		allowedAliases = []string{"COMMON"}
	}

	allowedAliasesSet := map[string]interface{}{}
	for _, a := range allowedAliases {
		allowedAliasesSet[strings.ToUpper(a)] = struct{}{}
	}

	uniqueAliases := UniqueAliases(str)

	count := 0
	for _, ua := range uniqueAliases {
		if _, ok := allowedAliasesSet[ua]; ok {
			continue
		}

		count++
	}
	return count > 1
}

func IsConfusable(str string, greedy bool, preferredAliases []string) []ConfusableResult {
	preferredAliasesSet := map[string]interface{}{}
	for _, a := range preferredAliases {
		preferredAliasesSet[strings.ToUpper(a)] = struct{}{}
	}

	outputs := []ConfusableResult{}
	checked := map[rune]interface{}{}
	for _, chr := range str {
		if _, ok := checked[chr]; ok {
			continue
		}
		checked[chr] = struct{}{}
		charAlias := Alias(chr)
		if _, ok := preferredAliasesSet[charAlias]; ok {
			// it's safe if the character might be confusable with homoglyphs from other
			// categories than our preferred categories (=aliases)
			continue
		}
		found, ok := confusablesData[string(chr)]
		if !ok {
			continue
		}
		// character λ is considered confusable if λ can be confused with a character from
		// preferred_aliases, e.g. if 'LATIN', 'ρ' is confusable with 'p' from LATIN.
		// if 'LATIN', 'Γ' is not confusable because in all the characters confusable with Γ,
		// none of them is LATIN.
		var potentiallyConfusable []Homoglyph
		if len(preferredAliasesSet) > 0 {
			potentiallyConfusable = []Homoglyph{}
		OUTER:
			for _, d := range found {
				for _, glyph := range d.C {
					a := Alias(glyph)
					if _, ok := preferredAliasesSet[a]; ok {
						potentiallyConfusable = found
						break OUTER
					}
				}
			}
		} else {
			potentiallyConfusable = found
		}

		if len(potentiallyConfusable) > 0 {
			outputs = append(outputs, ConfusableResult{
				Character:  chr,
				Alias:      charAlias,
				Homoglyphs: potentiallyConfusable,
			})
			if !greedy {
				return outputs
			}
		}
	}

	return outputs
}

func IsDangerous(str string, preferredAliases []string) bool {
	confusablesResult := IsConfusable(str, false, preferredAliases)
	return IsMixedScript(str, nil) && len(confusablesResult) > 0
}
