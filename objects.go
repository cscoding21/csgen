package csgen

import "strings"

// Struct a struct that abstracts a golang struct
type Struct struct {
	Name     string
	FilePath string
	Package  string
	Type     string
	Fields   []Field
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

// GetTag returns a single tag value by name based on the standard format rules
func (s *Field) GetTag(name string) string {
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
