package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/fredjn/etos/cli-go/pkg/logging"
)

// Client represents an ETOS client
type Client struct {
	Config     *config.Config
	HTTPClient *http.Client
	BaseURL    string
	Token      string
}

// NewClient creates a new ETOS client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		Config: cfg,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL: cfg.GetString("server.url"),
		Token:   cfg.GetString("server.token"),
	}
}

// DoRequest performs an HTTP request with the client's configuration
func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication token if present
	if c.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	// Add content type for POST requests
	if method == http.MethodPost && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	logging.DefaultLogger.Debug("Making request: %s %s", method, url)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// TestRun represents an ETOS test run
type TestRun struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetTestRun retrieves a test run by ID
func (c *Client) GetTestRun(ctx context.Context, id string) (*TestRun, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/testruns/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var testRunResp TestRunResponse
	if err := ParseResponse(resp, &testRunResp); err != nil {
		return nil, err
	}

	return testRunResp.TestRun, nil
}

// ListTestRuns retrieves a list of test runs
func (c *Client) ListTestRuns(ctx context.Context) ([]*TestRun, error) {
	resp, err := c.DoRequest(ctx, http.MethodGet, "/testruns", nil)
	if err != nil {
		return nil, err
	}

	var testRunsResp TestRunsResponse
	if err := ParseResponse(resp, &testRunsResp); err != nil {
		return nil, err
	}

	return testRunsResp.TestRuns, nil
} 