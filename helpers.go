package csgen

import (
	"fmt"

	"path/filepath"
	"strings"
)

func isPrimitive(t string) bool {
	et := getRawType(t)

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

func isSlice(t string) bool {
	return strings.Contains(t, "[]")
}

func isRefType(t string) bool {
	return strings.Contains(t, "*")
}

func getFileName(path string, name string) string {
	fullPath := filepath.Join(path, fmt.Sprintf("%s.gen.go", strings.ToLower(name)))
	return fullPath
}

func getFieldIndicator(source StructField, target StructField) string {
	if source.IsPointer == target.IsPointer {
		return ""
	}

	if source.IsPointer && !target.IsPointer {
		return "*"
	}

	return "&"
}

func getRawType(t string) string {
	out := t

	out = strings.ReplaceAll(out, "*", "")
	out = strings.ReplaceAll(out, "&", "")
	out = strings.ReplaceAll(out, "[]", "")

	return out
}

func stripPackageName(name string) string {
	elements := strings.Split(name, ".")

	if len(elements) > 1 {
		return elements[1]
	}
	return name
}

func extractPackageName(name string) string {
	elements := strings.Split(name, ".")

	if len(elements) > 1 {
		return elements[0]
	}
	return name
}

func isFullyQualifiedPackage(name string) bool {
	return strings.Contains(name, ".")
}

func sourceObjectContainsField(name string, graph Struct) bool {
	for _, f := range graph.Fields {
		if f.Name == name {
			return true
		}
	}

	return false
}

func objectSliceContainsName(name string, graph []Struct) *Struct {
	for _, o := range graph {
		if strings.ToLower(o.Name) == strings.ToLower(name) {
			return &o
		}
	}

	return nil
}

func getSourceField(name string, graph Struct) *StructField {
	for _, f := range graph.Fields {
		if f.Name == name {
			return &f
		}
	}

	return nil
}
