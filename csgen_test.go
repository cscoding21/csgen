package csgen

import (
	"fmt"
	"path/filepath"
	"strings"

	"os"
	"testing"
)

func TestGetStructs(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
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

func TestGetVars(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	vars, _ := GetVariables(fullPath)

	fmt.Println(vars)
}

func TestGetImports(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	imports, _ := GetImports(fullPath)
	expectedImports := 2

	fmt.Println(imports)

	//---this is fragile and should be rethought
	if len(imports) != expectedImports {
		t.Errorf("expected %v imports...got %v", expectedImports, len(imports))
	}
}

func TestGetInterfaces(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	interfaces, _ := GetInterfaces(fullPath)
	expectedImports := 3

	fmt.Println(interfaces)

	//---this is fragile and should be rethought
	if len(interfaces) != expectedImports {
		t.Errorf("expected %v interfaces...got %v", expectedImports, len(interfaces))
	}
}

func TestGetFields(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	structs, _ := GetStructs(fullPath)

	st := structs[0]

	nameField := st.GetField("Name")
	if nameField == nil {
		t.Error("expected field Name")
	}

	if nameField.Name != "Name" {
		t.Error("expected field Name")
	}
}

func TestGetFieldWithNil(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	structs, _ := GetStructs(fullPath)

	st := structs[0]

	noField := st.GetField("Zaphod")
	if noField != nil {
		t.Error("expected nil value")
	}
}

func TestGetTag(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "test_data.go")
	structs, _ := GetStructs(fullPath)

	testData := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "csval", want: "req"},
		{ok: true, have: "json", want: "id"},
		{ok: true, have: "zaphod", want: ""},
		{ok: true, have: "", want: ""},
	}

	st := structs[0]

	for _, input := range testData {
		field := st.GetField("ID")
		tag := field.GetTag(input.have)
		if tag != input.want {
			t.Errorf("for input %s, expected %s...got %s", input.have, input.want, tag)
		}
	}

	if st.GetField("Email").GetTag("csval") != "req,email" {
		t.Error("expected field ID to have tag csval:req")
	}
}

func TestGetFile(t *testing.T) {
	pwd, _ := os.Getwd()
	testFileName := "test_argument.go"
	filePath := filepath.Join(pwd, testFileName)
	file := GetFile(testFileName)

	//---test filename passed in by argument
	if file != filePath {
		t.Errorf("expected test_argument.go...got %s", file)
	}

	testFileName = "test_env.go"
	filePath = filepath.Join(pwd, testFileName)
	os.Setenv("GOFILE", testFileName)

	//---test filename in environment created by generator
	file = GetFile()
	if file != filePath {
		t.Errorf("expected test_env.go...got %s", file)
	}
}

func TestNewCSGenBuilderForFile(t *testing.T) {
	builder := NewCSGenBuilderForFile("test", "test")

	bs := builder.String()

	if !strings.Contains(bs, "package test") {
		t.Error("expected package test...got ", bs)
	}

	if !strings.Contains(bs, "// ### Common Sense Coding") {
		t.Error("expected import Common Sense Coding label ", bs)
	}
}

func TestExecuteTemplate(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	ts := TestStruct{Name: "test"}
	templateString := "test {{.Name}}"

	result := ExecuteTemplate("test", templateString, ts)

	if result != "test test" {
		t.Error("expected test test...got ", result)
	}
}

func TestWriteGeneratedGoFile(t *testing.T) {
	name := "test.go"
	contents := `package test
	
	func main() {}`

	err := WriteGeneratedGoFile(name, contents)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(name)

	bs, err := os.ReadFile(name)
	if err != nil {
		t.Error(err)
	}

	bsString := string(bs)
	if !strings.Contains(bsString, "package test") {
		t.Error("expected package test...got ", bsString)
	}

}
