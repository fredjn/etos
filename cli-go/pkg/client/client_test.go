package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://test.example.com",
		Token:   "test-token",
	}

	client := NewClient(cfg)
	assert.NotNil(t, client)
	assert.Equal(t, cfg.BaseURL, client.BaseURL)
	assert.Equal(t, cfg.Token, client.Token)
	assert.NotNil(t, client.HTTPClient)
}

func TestDoRequest(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
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
	client := NewClient(cfg)

	// Test request
	ctx := context.Background()
	resp, err := client.DoRequest(ctx, http.MethodGet, "/test", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestTestRunParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  TestRunParams
		wantErr bool
	}{
		{
			name: "valid params",
			params: TestRunParams{
				TestSuiteID: "test-suite",
				Identity:    "test-identity",
				Timeout:     time.Hour,
			},
			wantErr: false,
		},
		{
			name: "missing test suite",
			params: TestRunParams{
				Identity: "test-identity",
				Timeout:  time.Hour,
			},
			wantErr: true,
		},
		{
			name: "missing identity",
			params: TestRunParams{
				TestSuiteID: "test-suite",
				Timeout:     time.Hour,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			params: TestRunParams{
				TestSuiteID: "test-suite",
				Identity:    "test-identity",
				Timeout:     -time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
} 