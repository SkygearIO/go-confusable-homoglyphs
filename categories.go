package confusablehomoglyphs

//go:generate go run tools/gen.go

func AliasesCategories(chr rune) (string, string) {
	var l, r, c, m int
	l = 0
	r = len(categoryData.CodePointsRanges) - 1
	c = int(chr)

	// binary search
	for r >= l {
		m = (l + r) / 2
		if c < categoryData.CodePointsRanges[m][0] {
			r = m - 1
		} else if c > categoryData.CodePointsRanges[m][1] {
			l = m + 1
		} else {
			return categoryData.ISO15924aliases[categoryData.CodePointsRanges[m][2]],
				categoryData.Categories[categoryData.CodePointsRanges[m][3]]
		}
	}

	return "Unknown", "Zzzz"
}

func Alias(chr rune) string {
	a, _ := AliasesCategories(chr)
	return a
}

func Category(chr rune) string {
	_, c := AliasesCategories(chr)
	return c
}

func UniqueAliases(str string) []string {
	aSet := map[string]struct{}{}
	for _, chr := range str {
		aSet[Alias(chr)] = struct{}{}
	}

	keys := make([]string, len(aSet))
	i := 0
	for k := range aSet {
		keys[i] = k
		i++
	}

	return keys
}
