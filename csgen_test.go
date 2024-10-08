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
	fullPath := filepath.Join(pwd, "data_test.go")
	structs, err := GetStructs(fullPath)
	expectedCount := 3

	if err != nil {
		t.Error(err)
	}

	if len(structs) != expectedCount {
		t.Errorf("expected %v structs...got %v", expectedCount, len(structs))
	}

	for _, st := range structs {
		if st.Package != "csgen" {
			t.Errorf("struct %s: expected package csgen...got %v", st.Name, st.Package)
		}
	}
}

func TestGetVars(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	vars, _ := GetVariables(fullPath)

	fmt.Println(vars)
}

func TestGetImports(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	imports, _ := GetImports(fullPath)
	expectedImports := 2

	fmt.Println(imports)

	//---this is fragile and should be rethought
	if len(imports) != expectedImports {
		t.Errorf("expected %v imports...got %v", expectedImports, len(imports))
	}
}

func TestGetFunctions(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	functions, _ := GetFunctions(fullPath)
	expectedFunctions := 1

	fmt.Println(functions)

	//---this is fragile and should be rethought
	if len(functions) != expectedFunctions {
		t.Errorf("expected %v imports...got %v", expectedFunctions, len(functions))
	}

	for _, fn := range functions {
		if !fn.IsPublic {
			t.Errorf("expected function GetStructs...got %v", fn.Name)
		}
	}
}

func TestGetInterfaces(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
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
	fullPath := filepath.Join(pwd, "data_test.go")
	structs, _ := GetStructs(fullPath)

	st := structs[0]

	nameField := st.GetField("Name")
	if nameField == nil {
		t.Error("expected field Name")
		return
	}

	if !strings.EqualFold(nameField.Type, "string") {
		t.Error("expected field Name to be string")
	}
}

func TestGetFieldWithNil(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	structs, _ := GetStructs(fullPath)

	st := structs[0]

	noField := st.GetField("Zaphod")
	if noField != nil {
		t.Error("expected nil value")
	}
}

func TestGetTag(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
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

	testData := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "test_argument.go", want: filepath.Join(pwd, "test_argument.go")},
		{ok: true, have: "rel/test_argument.go", want: filepath.Join(pwd, "rel/test_argument.go")},
		{ok: true, have: "/home/test_argument.go", want: "/home/test_argument.go"},
	}

	//---evaluate argument test cases
	for _, input := range testData {
		f := GetFile(input.have)
		if !strings.EqualFold(f, input.want) {
			t.Errorf("for input %s, expected %s...got %s", input.have, input.want, f)
		}
	}

	//---test filename in environment created by generator
	testFileName := "test_env.go"
	filePath := filepath.Join(pwd, testFileName)
	os.Setenv("GOFILE", testFileName)
	file := GetFile()
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

	if !strings.Contains(bs, "// Code generated . DO NOT EDIT.") {
		t.Error("expected code generated indicator ", bs)
	}
}

func TestNewCSGenBuilderForOneOffFile(t *testing.T) {
	builder := NewCSGenBuilderForOneOffFile("test2", "test2")

	bs := builder.String()

	if !strings.Contains(bs, "package test2") {
		t.Error("expected package test2...got ", bs)
	}

	if !strings.Contains(bs, "// Developer Note: This file will only be generated once.") {
		t.Error("expected developer note", bs)
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

	import (
		"fmt"
	)
	
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

func TestMyFunc(t *testing.T) {
	//---this is a placeholder test to appease the code coverage tool.  The function is only used for testing
	MyFunc()
}

func TestProfileNode(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	f, err := GetAST(fullPath)
	if err != nil {
		t.Error(err)
	}

	for _, node := range f.Decls {
		ProfileNode(node)
	}
}

func TestGetStructFieldsWithEmbedded(t *testing.T) {
	pwd, _ := os.Getwd()
	fullPath := filepath.Join(pwd, "data_test.go")
	structs, _ := GetStructs(fullPath)

	st := structs[1]

	nameField := st.GetField("Name")
	if nameField == nil {
		t.Error("expected field Name")
		return
	}

	if !strings.EqualFold(nameField.Type, "string") {
		t.Error("expected field Name to be string")
	}
}

func TestGetCleanedType(t *testing.T) {
	testData := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "github.com/cscoding21/csgen/tmp.Location", want: "tmp.Location"},
	}

	for _, td := range testData {
		field := Field{
			Name:        "Field",
			Type:        td.have,
			IsPrimitive: false,
			IsPointer:   false,
			IsSlice:     false,
			IsPublic:    true,
		}

		ct := field.GetCleanedType()

		if !strings.EqualFold(ct, td.want) {
			t.Errorf("unexpected value: have %s - want %s", ct, td.want)
		}
	}
}

//----------------------------------------------------------------------------------

func TestLoadModule(t *testing.T) {
	cfg := GetDefaultPackageConfig()
	cfg.Tests = true
	module, err := LoadModule(cfg, "file=objects.go", "file=tests/tmp/tmp.go", "file=tests/common/cf.go")

	if err != nil {
		t.Error(err)
	}

	t.Log(module)
}
