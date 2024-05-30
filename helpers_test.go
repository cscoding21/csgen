package csgen

import (
	"testing"
)

func TestGetFileName(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "test.go", want: "here/test.gen.go"},
		{ok: true, have: "test", want: "here/test.gen.go"},
	}

	for _, input := range testCases {
		fn := getFileName("here", input.have)

		if fn != input.want {
			t.Errorf("expected %s...got %s", input.want, fn)
		}
	}
}

func TestGetFieldIndicator(t *testing.T) {
	testCases := []struct {
		ok   bool
		s1   Field
		s2   Field
		want string
	}{
		{ok: true, s1: Field{IsPointer: true}, s2: Field{IsPointer: true}, want: ""},
		{ok: true, s1: Field{IsPointer: true}, s2: Field{IsPointer: false}, want: "*"},
		{ok: true, s1: Field{IsPointer: false}, s2: Field{IsPointer: true}, want: "&"},
		{ok: true, s1: Field{IsPointer: false}, s2: Field{IsPointer: false}, want: ""},
	}

	for _, input := range testCases {
		fn := getFieldIndicator(input.s1, input.s2)

		if fn != input.want {
			t.Errorf("expected %s...got %s", input.want, fn)
		}
	}
}

func TestStripPackageName(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "Struct", want: "Struct"},
		{ok: true, have: "csgen.Struct", want: "Struct"},
		{ok: true, have: "csgen.StructField", want: "StructField"},
	}

	for _, input := range testCases {
		fn := stripPackageName(input.have)

		if fn != input.want {
			t.Errorf("expected %s...got %s", input.want, fn)
		}
	}
}

func TestExtractPackageName(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want string
	}{
		{ok: true, have: "Struct", want: ""},
		{ok: true, have: "csgen.Struct", want: "csgen"},
		{ok: true, have: "csgen.StructField", want: "csgen"},
	}

	for _, input := range testCases {
		fn := extractPackageName(input.have)

		if fn != input.want {
			t.Errorf("expected %s...got %s", input.want, fn)
		}
	}
}

func TestIsFullyQualifiedPackage(t *testing.T) {
	testCases := []struct {
		ok   bool
		have string
		want bool
	}{
		{ok: true, have: "Struct", want: false},
		{ok: true, have: "csgen.Struct", want: true},
		{ok: true, have: "csgen.StructField", want: true},
	}

	for _, input := range testCases {
		fn := isFullyQualifiedPackage(input.have)

		if fn != input.want {
			t.Errorf("expected %v...got %v", input.want, fn)
		}
	}
}

func TestSourceObjectContainsField(t *testing.T) {
	testStruct := Struct{
		Name: "TestStruct",
		Fields: []Field{
			{Name: "ID"},
			{Name: "Name"},
			{Name: "phone"},
		},
	}

	testCases := []struct {
		ok   bool
		st   Struct
		have string
		want bool
	}{
		{ok: true, st: testStruct, have: "ID", want: true},
		{ok: true, st: testStruct, have: "id", want: false},
		{ok: true, st: testStruct, have: "Name", want: true},
		{ok: true, st: testStruct, have: "Zaphod", want: false},
		{ok: true, st: testStruct, have: "phone", want: true},
	}

	for _, input := range testCases {
		fn := sourceObjectContainsField(input.have, input.st)

		if fn != input.want {
			t.Errorf("expected %v...got %v", input.want, fn)
		}
	}
}

func TestGetSructByName(t *testing.T) {
	slices := []Struct{
		{Name: "TestStruct"},
		{Name: "TestStruct2"},
		{Name: "TestStruct3"},
		{Name: "Test Struct3"},
	}

	testData := []struct {
		ok   bool
		have string
		want bool
	}{
		{ok: true, have: "TestStruct", want: true},
		{ok: true, have: "TestStruct2", want: true},
		{ok: true, have: "TestStruct3", want: true},
		{ok: true, have: "Test Struct3   ", want: false},
		{ok: true, have: "Test Struct3    ", want: true},
		{ok: true, have: "Zaphod", want: false},
	}

	for _, input := range testData {
		st := getStructByName(input.have, slices)
		if input.want && st == nil {
			t.Errorf("expected true for name %s...got false", input.have)
		}
	}
}

func TestIsPublic(t *testing.T) {
	testData := []struct {
		ok   bool
		have string
		want bool
	}{
		{ok: true, have: "TestStruct", want: true},
		{ok: true, have: "test", want: false},
		{ok: true, have: "_test", want: false},
		{ok: true, have: "", want: false},
		{ok: true, have: "TEST", want: true},
	}

	for _, input := range testData {
		if isPublic(input.have) != input.want {
			t.Errorf("expected true for name %s...got false", input.have)
		}
	}
}
