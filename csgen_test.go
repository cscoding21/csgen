package csgen

import (
	"path/filepath"

	"os"
	"testing"
	"time"
)

// ---Structs for testing
type TestStruct3 struct {
	ID      string `json:"id" csval:"req"`
	Name    string
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

func TestGetStructs(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "csgen_test.go")
	structs, err := GetStructs(fullPath)

	if err != nil {
		t.Error(err)
	}

	if len(structs) != 2 {
		t.Errorf("expected 2 structs...got %v", len(structs))
	}

	for _, st := range structs {
		if st.Package != "csgen" {
			t.Errorf("struct %s: expected package csgen...got %v", st.Name, st.Package)
		}
	}
}

func TestGetTag(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "csgen_test.go")
	structs, _ := GetStructs(fullPath)

	st := structs[0]

	for _, f := range st.Fields {

		t.Logf("%s - %s", f.Name, f.GetTag("csval"))
	}

}
