package test_results

import (
	"fmt"

	"github.com/eiffel-community/eiffelevents-sdk-go/editions/orizaba"
	"github.com/fredjn/etos/cli-go/pkg/client/v0/events"
)

// TestResults handles test results from an ETOS test run
type TestResults struct {
	events *events.Events
}

// NewTestResults creates a new test results handler
func NewTestResults(events *events.Events) *TestResults {
	return &TestResults{
		events: events,
	}
}

// hasFailed checks if any test suite has failed
func (t *TestResults) hasFailed() bool {
	for _, suite := range t.events.MainSuites {
		if suite.Finished == nil {
			continue
		}
		finished := suite.Finished.Get().(*orizaba.TestSuiteFinished)
		if finished.Data.Outcome.Conclusion != "SUCCESSFUL" {
			return true
		}
	}
	return false
}

// failMessages builds fail messages from main suites errors
func (t *TestResults) failMessages() []string {
	var messages []string
	for _, suite := range t.events.MainSuites {
		if suite.Finished == nil {
			continue
		}
		finished := suite.Finished.Get().(*orizaba.TestSuiteFinished)
		if finished.Data.Outcome.Conclusion != "SUCCESSFUL" {
			started := suite.Started.Get().(*orizaba.TestSuiteStarted)
			messages = append(messages, fmt.Sprintf("%s: %s", started.Data.Name, finished.Data.Outcome.Description))
		}
	}
	return messages
}

// GetResults gets the results from an ETOS test run
func (t *TestResults) GetResults() (bool, string, error) {
	// Check if we have all required events
	if t.events.TERCC == nil {
		return false, "", fmt.Errorf("no TERCC event found")
	}
	if t.events.Activity == nil || t.events.Activity.Triggered == nil {
		return false, "", fmt.Errorf("no activity triggered event found")
	}
	if t.events.Activity.Canceled != nil {
		return false, "", fmt.Errorf("activity was canceled")
	}
	if len(t.events.MainSuites) == 0 {
		return false, "", fmt.Errorf("no main test suites found")
	}
	for _, suite := range t.events.MainSuites {
		if suite.Finished == nil {
			started := suite.Started.Get().(*orizaba.TestSuiteStarted)
			return false, "", fmt.Errorf("test suite %s not finished", started.Meta.ID)
		}
	}

	// Check if any test suite failed
	if !t.hasFailed() {
		return true, "Test suite finished successfully.", nil
	}

	// Get fail messages
	messages := t.failMessages()
	if len(messages) == 1 {
		return false, messages[0], nil
	}
	if len(messages) > 0 {
		return false, messages[len(messages)-1], nil
	}

	return false, "Test case failures during test suite execution", nil
} 