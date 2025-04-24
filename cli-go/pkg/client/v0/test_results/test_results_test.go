package test_results

import (
	"testing"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/eiffel-community/eiffelevents-sdk-go/editions/orizaba"
	"github.com/stretchr/testify/assert"
)

func TestHasFailed(t *testing.T) {
	// Create test events
	successfulSuite := &orizaba.TestSuiteFinished{
		Data: orizaba.TestSuiteFinishedData{
			Outcome: orizaba.TestSuiteOutcome{
				Conclusion: "SUCCESSFUL",
			},
		},
	}
	successfulSuiteAny, _ := eiffelevents.NewAny(successfulSuite)

	failedSuite := &orizaba.TestSuiteFinished{
		Data: orizaba.TestSuiteFinishedData{
			Outcome: orizaba.TestSuiteOutcome{
				Conclusion: "FAILED",
			},
		},
	}
	failedSuiteAny, _ := eiffelevents.NewAny(failedSuite)

	tests := []struct {
		name     string
		suite    *eiffelevents.Any
		expected bool
	}{
		{
			name:     "Successful suite",
			suite:    successfulSuiteAny,
			expected: false,
		},
		{
			name:     "Failed suite",
			suite:    failedSuiteAny,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasFailed(tt.suite)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFailMessages(t *testing.T) {
	// Create test events
	suiteStarted := &orizaba.TestSuiteStarted{
		Data: orizaba.TestSuiteStartedData{
			Name:        "Test Suite",
			Description: "Test Description",
		},
	}
	suiteStartedAny, _ := eiffelevents.NewAny(suiteStarted)

	suiteFinished := &orizaba.TestSuiteFinished{
		Data: orizaba.TestSuiteFinishedData{
			Outcome: orizaba.TestSuiteOutcome{
				Conclusion:  "FAILED",
				Description: "Test failed",
			},
		},
	}
	suiteFinishedAny, _ := eiffelevents.NewAny(suiteFinished)

	tests := []struct {
		name     string
		suite    *TestSuite
		expected string
	}{
		{
			name: "Failed suite with description",
			suite: &TestSuite{
				Started:  suiteStartedAny,
				Finished: suiteFinishedAny,
			},
			expected: "Test Suite: Test failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := failMessages(tt.suite)
			assert.Equal(t, tt.expected, message)
		})
	}
}

func TestGetResults(t *testing.T) {
	// Create test events
	tercc := &orizaba.TestExecutionRecipeCollectionCreated{
		Meta: orizaba.Meta{
			ID: "test-tercc-id",
		},
	}
	terccAny, _ := eiffelevents.NewAny(tercc)

	activity := &orizaba.ActivityTriggered{
		Meta: orizaba.Meta{
			ID: "test-activity-id",
		},
	}
	activityAny, _ := eiffelevents.NewAny(activity)

	suiteStarted := &orizaba.TestSuiteStarted{
		Meta: orizaba.Meta{
			ID: "test-suite-id",
		},
		Data: orizaba.TestSuiteStartedData{
			Name: "Test Suite",
		},
	}
	suiteStartedAny, _ := eiffelevents.NewAny(suiteStarted)

	suiteFinished := &orizaba.TestSuiteFinished{
		Meta: orizaba.Meta{
			ID: "test-suite-id",
		},
		Data: orizaba.TestSuiteFinishedData{
			Outcome: orizaba.TestSuiteOutcome{
				Conclusion: "SUCCESSFUL",
			},
		},
	}
	suiteFinishedAny, _ := eiffelevents.NewAny(suiteFinished)

	tests := []struct {
		name           string
		events         *Events
		expectedResult *Results
		expectError    bool
	}{
		{
			name: "Successful test run",
			events: &Events{
				TERCC: terccAny,
				Activity: &Activity{
					Triggered: activityAny,
				},
				MainSuites: []*TestSuite{
					{
						Started:  suiteStartedAny,
						Finished: suiteFinishedAny,
					},
				},
			},
			expectedResult: &Results{
				Status: "SUCCESSFUL",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetResults(tt.events)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
} 