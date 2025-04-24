package events

import (
	"context"
	"testing"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/eiffel-community/eiffelevents-sdk-go/editions/orizaba"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	suiteResponse           *eiffelevents.Any
	activityResponse       *eiffelevents.Any
	canceledResponse      *eiffelevents.Any
	finishedResponse      *eiffelevents.Any
	mainSuitesResponse   []*eiffelevents.Any
	suiteFinishedResponse *eiffelevents.Any
}

func (m *mockRepository) RequestSuite(ctx context.Context, terccID string) (*eiffelevents.Any, error) {
	return m.suiteResponse, nil
}

func (m *mockRepository) RequestActivity(ctx context.Context, terccID string) (*eiffelevents.Any, error) {
	return m.activityResponse, nil
}

func (m *mockRepository) RequestActivityCanceled(ctx context.Context, activityID string) (*eiffelevents.Any, error) {
	return m.canceledResponse, nil
}

func (m *mockRepository) RequestActivityFinished(ctx context.Context, activityID string) (*eiffelevents.Any, error) {
	return m.finishedResponse, nil
}

func (m *mockRepository) RequestMainTestSuitesStarted(ctx context.Context, activityID string) ([]*eiffelevents.Any, error) {
	return m.mainSuitesResponse, nil
}

func (m *mockRepository) RequestTestSuiteFinished(ctx context.Context, suiteID string) (*eiffelevents.Any, error) {
	return m.suiteFinishedResponse, nil
}

func TestCollectActivity(t *testing.T) {
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

	canceled := &orizaba.ActivityCanceled{
		Meta: orizaba.Meta{
			ID: "test-canceled-id",
		},
	}
	canceledAny, _ := eiffelevents.NewAny(canceled)

	finished := &orizaba.ActivityFinished{
		Meta: orizaba.Meta{
			ID: "test-finished-id",
		},
	}
	finishedAny, _ := eiffelevents.NewAny(finished)

	tests := []struct {
		name           string
		repository     *mockRepository
		expectedEvents *Events
		expectError    bool
	}{
		{
			name: "Successful collection",
			repository: &mockRepository{
				suiteResponse:      terccAny,
				activityResponse:   activityAny,
				canceledResponse:   canceledAny,
				finishedResponse:   finishedAny,
			},
			expectedEvents: &Events{
				TERCC: terccAny,
				Activity: &Activity{
					Triggered: activityAny,
					Canceled:  canceledAny,
					Finished:  finishedAny,
				},
			},
			expectError: false,
		},
		{
			name: "Missing TERCC event",
			repository: &mockRepository{
				suiteResponse: nil,
			},
			expectedEvents: &Events{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector(nil, tt.repository)
			events, err := collector.CollectActivity(context.Background(), "test-id")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEvents, events)
			}
		})
	}
}

func TestCollect(t *testing.T) {
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
		repository     *mockRepository
		expectedEvents *Events
		expectError    bool
	}{
		{
			name: "Successful collection",
			repository: &mockRepository{
				suiteResponse:         terccAny,
				activityResponse:      activityAny,
				mainSuitesResponse:   []*eiffelevents.Any{suiteStartedAny},
				suiteFinishedResponse: suiteFinishedAny,
			},
			expectedEvents: &Events{
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
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector(nil, tt.repository)
			events, err := collector.Collect(context.Background(), "test-id")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedEvents, events)
			}
		})
	}
} 