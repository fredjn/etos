package v0

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/client"
	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://test.example.com",
		Token:   "test-token",
	}

	v0Client := NewClient(cfg)
	assert.NotNil(t, v0Client)
	assert.Equal(t, client.VersionV0, v0Client.Version())
}

func TestStartTestRun(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/testruns", r.URL.Path)

			var params map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&params)
			require.NoError(t, err)

			assert.Equal(t, "test-suite", params["test_suite_id"])
			assert.Equal(t, "test-identity", params["identity"])
			assert.Equal(t, float64(3600), params["timeout"])

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id": "test-id", "status": "running"}`))
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	// Create client
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Token:   "test-token",
	}
	v0Client := NewClient(cfg)

	// Test start test run
	ctx := context.Background()
	params := client.TestRunParams{
		TestSuiteID: "test-suite",
		Identity:    "test-identity",
		Timeout:     time.Hour,
	}

	testRun, err := v0Client.StartTestRun(ctx, params)
	require.NoError(t, err)
	assert.Equal(t, "test-id", testRun.ID)
	assert.Equal(t, "running", testRun.Status)
}

func TestStopTestRun(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/testruns/test-id/stop", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	// Create client
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Token:   "test-token",
	}
	v0Client := NewClient(cfg)

	// Test stop test run
	ctx := context.Background()
	err := v0Client.StopTestRun(ctx, "test-id")
	require.NoError(t, err)
}

func TestGetTestRun(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/testruns/test-id", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "test-id", "status": "running"}`))
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	// Create client
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Token:   "test-token",
	}
	v0Client := NewClient(cfg)

	// Test get test run
	ctx := context.Background()
	testRun, err := v0Client.GetTestRun(ctx, "test-id")
	require.NoError(t, err)
	assert.Equal(t, "test-id", testRun.ID)
	assert.Equal(t, "running", testRun.Status)
}

func TestListTestRuns(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/testruns", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "test-id-1", "status": "running"}, {"id": "test-id-2", "status": "completed"}]`))
		}),
	}
	go server.ListenAndServe()
	defer server.Close()

	// Create client
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Token:   "test-token",
	}
	v0Client := NewClient(cfg)

	// Test list test runs
	ctx := context.Background()
	testRuns, err := v0Client.ListTestRuns(ctx)
	require.NoError(t, err)
	assert.Len(t, testRuns, 2)
	assert.Equal(t, "test-id-1", testRuns[0].ID)
	assert.Equal(t, "running", testRuns[0].Status)
	assert.Equal(t, "test-id-2", testRuns[1].ID)
	assert.Equal(t, "completed", testRuns[1].Status)
} 