package csgen

import (
	"fmt"
	"regexp"
	"strings"
)

// Module represents an object graph for the entire codebase
type Module struct {
	Name     string
	Path     string
	Packages []Package
}

// Package represents an object graph for an entire package
type Package struct {
	ID       string
	Name     string
	FullName string
	Path     string
	Files    []string
	Structs  []Struct
}

// Struct a struct that abstracts a golang struct
type Struct struct {
	Name           string
	FilePath       string
	Package        string
	Type           string
	Fields         []Field
	EmbeddedFields []Field
}

// Field a struct that represents a single field within a struct abstraction
type Field struct {
	Name        string
	Type        string
	TagString   string
	IsPrimitive bool
	IsPointer   bool
	IsSlice     bool
	IsPublic    bool
}

// Function a struct that represents a single function abstraction
type Function struct {
	Name      string
	Receiver  *string
	Arguments []Field
	Returns   []Field
	IsPublic  bool
}

// Interface a struct that represents a single interface abstraction
type Interface struct {
	Name     string
	Methods  []Function
	IsPublic bool
}

// GetPackage return a package from a module based on its name
func (m *Module) GetPackage(name string) *Package {
	for i, p := range m.Packages {
		if strings.EqualFold(p.Name, name) || strings.EqualFold(p.ID, name) {
			return &m.Packages[i]
		}
	}

	return nil
}

// GetStruct return a struct with the name that matches the argument
func (p *Package) GetStruct(name string) *Struct {
	for i, s := range p.Structs {
		if strings.EqualFold(s.Name, name) {
			return &p.Structs[i]
		}
	}

	return nil
}

// GetTag returns a single tag value by name based on the standard format rules
func (f *Field) GetTag(name string) string {
	if len(f.TagString) == 0 || len(name) == 0 {
		return ""
	}

	tagString := strings.Trim(f.TagString, "`")
	if strings.HasPrefix(tagString, "\"") {
		tagString = strings.Trim(tagString, "\"")
	}

	tags := strings.Split(tagString, " ")
	for _, t := range tags {
		if strings.HasPrefix(t, name) {
			if strings.Contains(t, ":") {
				out := strings.Split(t, ":")[1]
				if strings.HasPrefix(out, "\"") {
					out = strings.Trim(out, "\"")
				}

				return out
			}

			return t
		}
	}

	return ""
}

// GetCleanedType return a cleaned version of the type in case it includes a fully qualified namespace
func (f *Field) GetCleanedType() string {
	if strings.Contains(f.Type, "/") {
		re := regexp.MustCompile(`[\w-\.]`)
		ta := strings.Split(f.Type, "/")

		nakedType := ta[len(ta)-1]

		indicator := re.ReplaceAllString(ta[0], "")

		return fmt.Sprintf("%s%s", indicator, nakedType)
	}

	return f.Type
}

// GetField return a field object of a struct by its name
func (s *Struct) GetField(name string) *Field {
	for _, f := range s.Fields {
		if f.Name == name {
			return &f
		}
	}

	return nil
}

// ContainsField returns true if the struct contains a field with the passed in name
func (s *Struct) ContainsField(name string) bool {
	f := s.GetField(name)

	return f != nil
}

// AllFields return all fields and embedded fields
func (s *Struct) AllFields() []Field {
	return append(s.Fields, s.EmbeddedFields...)
}
