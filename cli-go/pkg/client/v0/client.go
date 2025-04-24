package v0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/client"
	"github.com/fredjn/etos/cli-go/pkg/client/v0/events"
	"github.com/fredjn/etos/cli-go/pkg/client/v0/test_results"
	"github.com/fredjn/etos/cli-go/pkg/config"
)

// Client represents an ETOS v0 client
type Client struct {
	*client.Client
}

// NewClient creates a new ETOS v0 client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		Client: client.NewClient(cfg),
	}
}

// Version returns the API version
func (c *Client) Version() client.Version {
	return client.VersionV0
}

// StartTestRun starts a new test run
func (c *Client) StartTestRun(ctx context.Context, params client.TestRunParams) (*client.TestRun, error) {
	// Check environment variables for missing parameters
	if params.Identity == "" {
		params.Identity = os.Getenv("IDENTITY")
	}
	if params.TestSuiteID == "" {
		params.TestSuiteID = os.Getenv("TEST_SUITE")
	}
	if params.ParentActivityID == "" {
		if actt := os.Getenv("EIFFEL_ACTIVITY_TRIGGERED"); actt != "" {
			var activity struct {
				Meta struct {
					ID string `json:"id"`
				} `json:"meta"`
			}
			if err := json.Unmarshal([]byte(actt), &activity); err == nil {
				params.ParentActivityID = activity.Meta.ID
			}
		}
	}

	// Convert parameters to v0 format
	v0Params := map[string]interface{}{
		"test_suite_id": params.TestSuiteID,
		"identity":      params.Identity,
		"environment":   params.Environment,
		"timeout":       int(params.Timeout.Seconds()),
	}

	if params.ParentActivityID != "" {
		v0Params["parent_activity_id"] = params.ParentActivityID
	}
	if params.Workspace != "" {
		v0Params["workspace"] = params.Workspace
	}
	if params.ArtifactDir != "" {
		v0Params["artifact_dir"] = params.ArtifactDir
	}
	if params.ReportDir != "" {
		v0Params["report_dir"] = params.ReportDir
	}
	if params.ProviderConfig.IUTProvider != "" {
		v0Params["iut_provider"] = params.ProviderConfig.IUTProvider
	}
	if params.ProviderConfig.ExecutionSpaceProvider != "" {
		v0Params["execution_space_provider"] = params.ProviderConfig.ExecutionSpaceProvider
	}
	if params.ProviderConfig.LogAreaProvider != "" {
		v0Params["log_area_provider"] = params.ProviderConfig.LogAreaProvider
	}
	if len(params.Dataset) > 0 {
		v0Params["dataset"] = params.Dataset
	}

	// Make request
	resp, err := c.DoRequest(ctx, http.MethodPost, "/testruns", v0Params)
	if err != nil {
		return nil, fmt.Errorf("failed to start test run: %w", err)
	}
	defer resp.Body.Close()

	var testRun client.TestRun
	if err := json.NewDecoder(resp.Body).Decode(&testRun); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &testRun, nil
}

// StopTestRun stops a running test run
func (c *Client) StopTestRun(ctx context.Context, id string) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/testruns/%s/stop", id), nil)
	if err != nil {
		return fmt.Errorf("failed to stop test run: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// GetTestRun retrieves a test run by ID
func (c *Client) GetTestRun(ctx context.Context, id string) (*client.TestRun, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/testruns/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get test run: %w", err)
	}
	defer resp.Body.Close()

	var testRun client.TestRun
	if err := json.NewDecoder(resp.Body).Decode(&testRun); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &testRun, nil
}

// ListTestRuns retrieves a list of test runs
func (c *Client) ListTestRuns(ctx context.Context) ([]*client.TestRun, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/testruns", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list test runs: %w", err)
	}
	defer resp.Body.Close()

	var testRuns []*client.TestRun
	if err := json.NewDecoder(resp.Body).Decode(&testRuns); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return testRuns, nil
}

// GetTestResults gets the test results for a test run
func (c *Client) GetTestResults(ctx context.Context, terccID string, repository events.EventRepository) (bool, string, error) {
	collector := events.NewCollector(c, repository)
	events, err := collector.Collect(ctx, terccID)
	if err != nil {
		return false, "", fmt.Errorf("failed to collect events: %w", err)
	}

	results := test_results.NewTestResults(events)
	return results.GetResults()
} 