package events

import (
	"context"
	"fmt"
	"time"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/eiffel-community/eiffelevents-sdk-go/editions/orizaba"
	"github.com/fredjn/etos/cli-go/pkg/client"
	"github.com/fredjn/etos/cli-go/pkg/client/v0"
)

// Collector collects events from an event repository
type Collector struct {
	client     *v0.Client
	repository EventRepository
	events     *Events
}

// NewCollector creates a new event collector
func NewCollector(client *v0.Client, repository EventRepository) *Collector {
	return &Collector{
		client:     client,
		repository: repository,
		events:     &Events{},
	}
}

// EventRepository defines the interface for event repositories
type EventRepository interface {
	RequestSuite(ctx context.Context, terccID string) (*eiffelevents.Any, error)
	RequestActivity(ctx context.Context, terccID string) (*eiffelevents.Any, error)
	RequestActivityCanceled(ctx context.Context, activityID string) (*eiffelevents.Any, error)
	RequestActivityFinished(ctx context.Context, activityID string) (*eiffelevents.Any, error)
	RequestMainTestSuitesStarted(ctx context.Context, activityID string) ([]*eiffelevents.Any, error)
	RequestTestSuiteFinished(ctx context.Context, suiteID string) (*eiffelevents.Any, error)
}

// CollectActivity collects activity events from ETOS
func (c *Collector) CollectActivity(ctx context.Context, terccID string) (*Events, error) {
	// Get TERCC event
	if c.events.TERCC == nil {
		tercc, err := c.repository.RequestSuite(ctx, terccID)
		if err != nil {
			return nil, fmt.Errorf("failed to get TERCC event: %w", err)
		}
		c.events.TERCC = tercc
	}
	if c.events.TERCC == nil {
		return c.events, nil
	}

	// Get activity events
	if c.events.Activity == nil {
		c.events.Activity = &Activity{}
	}

	// Get triggered event
	if c.events.Activity.Triggered == nil {
		triggered, err := c.repository.RequestActivity(ctx, terccID)
		if err != nil {
			return nil, fmt.Errorf("failed to get activity triggered event: %w", err)
		}
		c.events.Activity.Triggered = triggered
	}
	if c.events.Activity.Triggered == nil {
		return c.events, nil
	}

	// Get canceled event
	if c.events.Activity.Canceled == nil {
		canceled, err := c.repository.RequestActivityCanceled(ctx, c.events.Activity.Triggered.Get().(*orizaba.ActivityTriggered).Meta.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get activity canceled event: %w", err)
		}
		c.events.Activity.Canceled = canceled
	}

	// Get finished event
	if c.events.Activity.Finished == nil {
		finished, err := c.repository.RequestActivityFinished(ctx, c.events.Activity.Triggered.Get().(*orizaba.ActivityTriggered).Meta.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get activity finished event: %w", err)
		}
		c.events.Activity.Finished = finished
	}

	return c.events, nil
}

// Collect collects all events from ETOS
func (c *Collector) Collect(ctx context.Context, terccID string) (*Events, error) {
	if c.events.Activity == nil || c.events.Activity.Triggered == nil {
		events, err := c.CollectActivity(ctx, terccID)
		if err != nil {
			return nil, err
		}
		if events.Activity == nil || events.Activity.Triggered == nil {
			return c.events, nil
		}
	}

	// Get main test suites
	started, err := c.repository.RequestMainTestSuitesStarted(ctx, c.events.Activity.Triggered.Get().(*orizaba.ActivityTriggered).Meta.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get main test suites: %w", err)
	}

	// Create test suites
	suites := make(map[string]*TestSuite)
	for _, s := range c.events.MainSuites {
		if s.Started != nil {
			suites[s.Started.Get().(*orizaba.TestSuiteStarted).Meta.ID] = s
		}
	}

	// Add new test suites
	for _, s := range started {
		if _, exists := suites[s.Get().(*orizaba.TestSuiteStarted).Meta.ID]; !exists {
			suites[s.Get().(*orizaba.TestSuiteStarted).Meta.ID] = &TestSuite{
				Started: s,
			}
		}
	}

	// Get finished events for all suites
	for id, suite := range suites {
		if suite.Finished == nil {
			finished, err := c.repository.RequestTestSuiteFinished(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("failed to get test suite finished event: %w", err)
			}
			suite.Finished = finished
		}
	}

	// Convert suites map to slice
	c.events.MainSuites = make([]*TestSuite, 0, len(suites))
	for _, suite := range suites {
		c.events.MainSuites = append(c.events.MainSuites, suite)
	}

	return c.events, nil
} 