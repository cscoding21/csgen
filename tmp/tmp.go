package tmp

import (
	"fmt"
)

// DoAThing a function for testing
func DoAThing() {
	fmt.Println("doing...")
}

// DoType a struct for testing
type DoType struct {
	Do       bool
	What     string
	Location Location

	EmbedMe
}
