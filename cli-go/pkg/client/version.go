package client

import (
	"context"
	"time"
)

// Version represents an ETOS API version
type Version string

const (
	// VersionV0 represents ETOS API version 0
	VersionV0 Version = "v0"
	// VersionV1Alpha represents ETOS API version 1alpha
	VersionV1Alpha Version = "v1alpha"
)

// ProviderConfig represents configuration for different providers
type ProviderConfig struct {
	IUTProvider            string
	ExecutionSpaceProvider string
	LogAreaProvider        string
}

// VersionedClient represents a version-specific ETOS client
type VersionedClient interface {
	// Version returns the API version
	Version() Version

	// StartTestRun starts a new test run
	StartTestRun(ctx context.Context, params TestRunParams) (*TestRun, error)

	// StopTestRun stops a running test run
	StopTestRun(ctx context.Context, id string) error

	// GetTestRun retrieves a test run by ID
	GetTestRun(ctx context.Context, id string) (*TestRun, error)

	// ListTestRuns retrieves a list of test runs
	ListTestRuns(ctx context.Context) ([]*TestRun, error)
}

// TestRunParams represents parameters for starting a test run
type TestRunParams struct {
	// TestSuiteID is the ID of the test suite to run
	TestSuiteID string `json:"test_suite_id"`
	// Identity is the artifact created identity purl or ID
	Identity string `json:"identity"`
	// ParentActivityID is the activity for the TERCC to link to
	ParentActivityID string `json:"parent_activity_id"`
	// Workspace is the workspace to do all the work in
	Workspace string `json:"workspace"`
	// ArtifactDir is where test artifacts should be stored
	ArtifactDir string `json:"artifact_dir"`
	// ReportDir is where test reports should be stored
	ReportDir string `json:"report_dir"`
	// Environment is the environment configuration
	Environment map[string]interface{} `json:"environment"`
	// Timeout is the maximum duration for the test run
	Timeout time.Duration `json:"timeout"`
	// ProviderConfig contains provider-specific configurations
	ProviderConfig ProviderConfig `json:"provider_config"`
	// Dataset contains additional dataset information for the environment provider
	Dataset []map[string]interface{} `json:"dataset"`
} 