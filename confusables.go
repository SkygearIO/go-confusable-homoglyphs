package confusablehomoglyphs

import (
	"bytes"
	"strings"
)

type ConfusableResult struct {
	Character  rune        `json:"character"`
	Alias      string      `json:"alias"`
	Homoglyphs []Homoglyph `json:"homoglyphs"`
}

// IsMixedScript checks if str contains mixed-scripts content,
// excluding script blocks aliases in allowedAliases.
// E.g. ``B. C`` is not considered mixed-scripts by default: it contains characters
// from Latin and Common, but Common is excluded by default.
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

// IsConfusable check if str contains characters which might be confusable with
// characters from preferredAliases.
// If greedy is false, it will only return the first confusable character
// found without looking at the rest of the string, greedy is true returns
// all of them.
// preferredAliases can take an array of unicode block aliases to
// be considered as your 'base' unicode blocks
func IsConfusable(str string, greedy bool, preferredAliases []string) []ConfusableResult {
	preferredAliasesSet := map[string]struct{}{}
	for _, a := range preferredAliases {
		preferredAliasesSet[strings.ToUpper(a)] = struct{}{}
	}

	outputs := []ConfusableResult{}
	checked := map[rune]struct{}{}
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

// SetConfusableToLatin check if str contains characters which might be confusable with
// characters from preferredAliases and normalizes it to latin character.
// it will return all of the confusable characters converted to LATIN.
func SetConfusableToLatin(str string, preferredAliases []string) string {
	preferredAliasesSet := map[string]struct{}{}

	for _, a := range preferredAliases {
		preferredAliasesSet[strings.ToUpper(a)] = struct{}{}
	}

	var outputs bytes.Buffer
	checked := map[rune]string{}

	for _, chr := range str {
		// if character is already checked and if its converted to LATIN, then we use the LATIN value
		// instead of converting again
		if _, ok := checked[chr]; ok {
			convertedVal := checked[chr]
			if convertedVal == "" {
				outputs.WriteString(string(chr))
				continue
			} else {
				outputs.WriteString(string(convertedVal))
				continue
			}
		}

		checked[chr] = ""
		charAlias := Alias(chr)

		for _, expectedAlias := range preferredAliases {
			if charAlias != strings.ToUpper(expectedAlias) {

				found, ok := confusablesData[string(chr)]
				if !ok {
					continue
				}

				var potentiallyConfusable []Homoglyph
				if len(preferredAliasesSet) > 0 {
					potentiallyConfusable = []Homoglyph{}
				OUTER:
					for _, d := range found {
						for _, glyph := range d.C {
							latinChar := d.C
							a := Alias(glyph)

							if res := strings.EqualFold(a, "LATIN"); res {
								if _, ok := preferredAliasesSet[a]; ok {
									potentiallyConfusable = found
									outputs.WriteString(string(latinChar))
									checked[chr] = latinChar
									break OUTER
								}
							}
						}
					}
				} else {
					potentiallyConfusable = found
				}

				if len(potentiallyConfusable) > 0 {
				}
			}
		}

		if _, ok := preferredAliasesSet[charAlias]; ok {
			outputs.WriteString(string(chr))
			continue
		}
	}
	return outputs.String()
}

// IsDangerous checks if str can be dangerous, i.e. is it not only mixed-scripts
// but also contains characters from other scripts than the ones in preferredAliases
// that might be confusable with characters from scripts in preferredAliases.
func IsDangerous(str string, preferredAliases []string) bool {
	confusablesResult := IsConfusable(str, false, preferredAliases)
	return IsMixedScript(str, nil) && len(confusablesResult) > 0
}
