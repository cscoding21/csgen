package csgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"unicode"
	"unicode/utf8"

	"path/filepath"
	"strings"
)

// GetAST return an AST object from a single file
func GetAST(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, filePath, nil, 0)
}

// IsPrimitive return true if the type is a primitive
func IsPrimitive(t string) bool {
	et := GetRawType(t)

	return et == "string" ||
		et == "int" ||
		et == "int8" ||
		et == "int16" ||
		et == "int32" ||
		et == "int64" ||
		et == "uint" ||
		et == "uint8" ||
		et == "uint16" ||
		et == "uint32" ||
		et == "uint64" ||
		et == "float32" ||
		et == "float64" ||
		et == "bool" ||
		et == "time.Time" ||
		et == "byte"
}

// IsSlice return true if the type is a slice
func IsSlice(t string) bool {
	return strings.Contains(t, "[]")
}

// IsRefType return true if the type is a pointer
func IsRefType(t string) bool {
	return strings.Contains(t, "*")
}

// IsPublic return true if the type is public
func IsPublic(t string) bool {
	if len(t) == 0 {
		return false
	}

	// Get the first rune (character) of the string
	r, _ := utf8.DecodeRuneInString(t)

	// Check if the rune is an uppercase letter
	return unicode.IsUpper(r)
}

// GetFileName return a file name suffixed with ".gen.go" to indicate that is was generated
func GetFileName(imp string, path string, name string) string {
	//---remove the extension so it can be readded without issue
	name = strings.TrimSuffix(name, filepath.Ext(name))

	fullPath := filepath.Join(path, fmt.Sprintf("z_%s_%s.gen.go", strings.ToLower(name), strings.ToLower(imp)))
	return fullPath
}

// GetFieldIndicator for creating an assignment operations, returns an indicator for the field based on the type of the source and target
func GetFieldIndicator(source Field, target Field) string {
	if source.IsPointer == target.IsPointer {
		return ""
	}

	if source.IsPointer && !target.IsPointer {
		return "*"
	}

	return "&"
}

// GetRawType return the raw type of a type, removing the pointer, reference and slice indicators
func GetRawType(t string) string {
	out := t

	out = strings.ReplaceAll(out, "*", "")
	out = strings.ReplaceAll(out, "&", "")
	out = strings.ReplaceAll(out, "[]", "")

	return out
}

// StripPackageName removes the package name from a fully qualified name
func StripPackageName(name string) string {
	elements := strings.Split(name, ".")

	if len(elements) > 1 {
		return elements[1]
	}

	return name
}

// ExtractPackageName return the package name from a fully qualified name
func ExtractPackageName(name string) string {
	elements := strings.Split(name, ".")

	if len(elements) > 1 {
		return elements[0]
	}

	return ""
}

// IsFullyQualifiedPackage return true if the package name is fully qualified
func IsFullyQualifiedPackage(name string) bool {
	return strings.Contains(name, ".")
}

// GetStructByName given a slice of structs, return one if it matches the name
func GetStructByName(name string, graph []Struct) *Struct {
	for _, o := range graph {
		if strings.EqualFold(o.Name, strings.Trim(name, " ")) {
			return &o
		}
	}

	return nil
}
