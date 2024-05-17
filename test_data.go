package csgen

import (
	"fmt"
	"time"
)

var MyString string
var MyInt int

type MyInterface interface {
	Fizz()
}

type MyInterface2 interface {
	Buzz(input string) (string, error)
}

type MyInterface3 interface {
	FizzBuzz() (string, error)
}

// ---Structs for testing
type TestStruct3 struct {
	ID      string `json:"id" csval:"req"`
	Name    string
	Email   string `json:"email" csval:"req,email"`
	Number  int
	Boolean bool
	Time    time.Time
}

type TestStruct4 struct {
	ID      string "json:\"id\" csval:\"req\""
	Name    string
	Number  int
	Boolean bool
	Time    time.Time
}

func MyFunc() {
	fmt.Println("Hello World")
}
