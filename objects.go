package csgen

import "strings"

type Struct struct {
	Name     string
	FilePath string
	Package  string
	Type     string
	Fields   []StructField
}

type StructField struct {
	Name        string
	Type        string
	TagString   string
	IsPrimitive bool
	IsPointer   bool
	IsSlice     bool
}

// GetTag returns a single tag value by name based on the standard format rules
func (s StructField) GetTag(name string) string {
	if len(s.TagString) == 0 {
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
