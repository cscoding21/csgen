package tmp

import "time"

// EmbedMe a struct for testing embedded functionality
type EmbedMe struct {
	Field1 time.Time
	Field2 string
	Field3 float32
}

type Location struct {
	Lat float32
	Lon float32
}
