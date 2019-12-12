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
}
```
