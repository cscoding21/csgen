package csgen

import (
	"fmt"
	"time"
)

// ---Variables for testing

// MyString is a placeholder string to be identified by the tests
var MyString string

// MyInt is a placeholder int to be identified by the tests
var MyInt int

// MyInterface is a placeholder interface to be identified by the tests
type MyInterface interface {
	Fizz()
}

// MyInterface2 is a placeholder interface to be identified by the tests that contains arguments and returns
type MyInterface2 interface {
	Buzz(input string) (string, error)
}

// MyInterface3 is a placeholder interface to be identified by the tests that contains returns but no arguments
type MyInterface3 interface {
	FizzBuzz() (string, error)
}

// ---Structs for testing

// TestStruct3 is a placeholder struct to be identified and characterized by the tests
type TestStruct3 struct {
	ID      string `json:"id" csval:"req"`
	Name    string
	Email   string `json:"email" csval:"req,email"`
	Number  int
	Boolean bool
	Time    time.Time
}

// TestStruct4 is a placeholder struct to be identified and characterized by the tests
type TestStruct4 struct {
	ID      string "json:\"id\" csval:\"req\""
	Name    string
	Number  int
	Boolean bool
	Time    time.Time
}

// MyFunc is a placeholder function to be identified and characterized by the tests
func MyFunc() {
	fmt.Println("Hello World")
}
