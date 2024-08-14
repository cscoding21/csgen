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

// TestKey used for testing
type TestKey string

// Activity used for testing
type Activity struct {
	ID       string    `json:"id"`
	Type     string    `json:"type"`
	Summary  string    `json:"summary"`
	Detail   *string   `json:"detail,omitempty"`
	Context  string    `json:"context"`
	TargetID *string   `json:"targetID,omitempty"`
	Time     time.Time `json:"time"`
	Key      TestKey   `json:"key"`
}

// ActivityResults used for testing
type ActivityResults struct {
	Results []*Activity `json:"results,omitempty"`
}

// TestSource used for testing
type TestSource struct {
	ID       string
	Name     string
	Age      *int
	Location LocationSource
}

// LocationSource used for testing
type LocationSource struct {
	Lat float64
	Lon float64
}
