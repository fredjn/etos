package events

import (
	"github.com/eiffel-community/eiffelevents-sdk-go"
)

// Events represents a collection of Eiffel events
type Events struct {
	TERCC      *eiffelevents.Any
	Activity   *Activity
	MainSuites []*TestSuite
}

// Activity represents activity-related events
type Activity struct {
	Triggered *eiffelevents.Any
	Canceled  *eiffelevents.Any
	Finished  *eiffelevents.Any
}

// TestSuite represents test suite-related events
type TestSuite struct {
	Started  *eiffelevents.Any
	Finished *eiffelevents.Any
} 