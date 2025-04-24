package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse represents an error response from the ETOS server
type ErrorResponse struct {
	Error string `json:"error"`
}

// TestRunResponse represents a testrun response from the ETOS server
type TestRunResponse struct {
	TestRun *TestRun `json:"testrun"`
}

// TestRunsResponse represents a list of testruns response from the ETOS server
type TestRunsResponse struct {
	TestRuns []*TestRun `json:"testruns"`
}

// parseResponse parses the HTTP response body into the target interface
func parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("failed to parse error response: %w", err)
		}
		return fmt.Errorf("server error: %s", errResp.Error)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
} 