package events

import "time"

// Activity represents ETOS activity events
type Activity struct {
	Triggered *Event `json:"triggered,omitempty"`
	Canceled  *Event `json:"canceled,omitempty"`
	Finished  *Event `json:"finished,omitempty"`
}

// SubSuite represents ETOS sub suite events
type SubSuite struct {
	Started  *Event `json:"started,omitempty"`
	Finished *Event `json:"finished,omitempty"`
}

// TestSuite represents ETOS main suite events
type TestSuite struct {
	Started    *Event      `json:"started,omitempty"`
	Finished   *Event      `json:"finished,omitempty"`
	SubSuites  []*SubSuite `json:"sub_suites,omitempty"`
}

// Events represents all ETOS events
type Events struct {
	Activity   *Activity    `json:"activity,omitempty"`
	TERCC      *Event       `json:"tercc,omitempty"`
	MainSuites []*TestSuite `json:"main_suites,omitempty"`
}

// Event represents a generic ETOS event
type Event struct {
	Meta struct {
		ID        string    `json:"id"`
		Type      string    `json:"type"`
		Version   string    `json:"version"`
		Time      time.Time `json:"time"`
		Tags      []string  `json:"tags,omitempty"`
		Source    string    `json:"source"`
		DomainID  string    `json:"domainId"`
		SchemaURI string    `json:"schemaUri"`
	} `json:"meta"`
	Data map[string]interface{} `json:"data"`
}

// TestSuiteOutcome represents the outcome of a test suite
type TestSuiteOutcome struct {
	Verdict     string `json:"verdict"`
	Conclusion  string `json:"conclusion"`
	Description string `json:"description"`
} 