package csgen

import "strings"

// Struct a struct that abstracts a golang struct
type Struct struct {
	Name     string
	FilePath string
	Package  string
	Type     string
	Fields   []StructField
}

// StructField a struct that represents a single field within a struct abstraction
type StructField struct {
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
	Arguments []StructField
	Returns   []StructField
	IsPublic  bool
}

// Interface a struct that represents a single interface abstraction
type Interface struct {
	Name     string
	Methods  []Function
	IsPublic bool
}

// GetTag returns a single tag value by name based on the standard format rules
func (s *StructField) GetTag(name string) string {
	if len(s.TagString) == 0 || len(name) == 0 {
		return ""
	}

	tagString := strings.Trim(s.TagString, "`")
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

// GetField return a field object of a struct by its name
func (s *Struct) GetField(name string) *StructField {
	for _, f := range s.Fields {
		if f.Name == name {
			return &f
		}
	}

	return nil
}
