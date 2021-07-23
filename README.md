Forked
Forked from https://github.com/SkygearIO/go-confusable-homoglyphs on 20July2021 at this commit to support conversion of confusable characters to LATIN characters. We need this to support normalization of domains such as 'ρhishlabs' which include confusable characters such as 'ρ' which is converted to LATIN 'p' 'phishlabs'.

SkygearIO (the author) is currently supporting to check if the characters are confusable characters within any character set.

# confusable_homoglyphs

Golang version of Python library https://github.com/vhf/confusable_homoglyphs

# Example

```go
package main

import (
	"fmt"

	confusable "github.com/skygeario/go-confusable-homoglyphs"
)

func main() {

	// detact mixed script and confusable
	isDangerous := confusable.IsDangerous("AlaskaJazz", []string{})
	fmt.Println(isDangerous) // should be false

	isDangerous = confusable.IsDangerous("ΑlaskaJazz", []string{})
	fmt.Println(isDangerous) // should be true

	// detect confusable with preferred aliases
	confusables := confusable.IsConfusable("microsoft", false, []string{"latin", "common"})
	fmt.Println(confusables) // should be empty

	confusables = confusable.IsConfusable("microsоft", false, []string{"latin", "common"})
	fmt.Println(confusables) // should show confusable homoglyphs

	// detect confusable with preferred aliases
	latinResult := confusable.SetConfusableToLatin("ρhishlabs", []string{"latin"})
	fmt.Println(latinResult) // should show converted latin string result - phishlabs

	latinResult = confusable.SetConfusableToLatin("ᑲankofamericα", []string{"latin"})
	fmt.Println(latinResult) // should show converted latin string result - bankofamerica
}
```
