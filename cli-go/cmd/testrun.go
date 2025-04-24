package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/fredjn/etos/cli-go/pkg/client"
	"github.com/fredjn/etos/cli-go/pkg/client/v0"
	"github.com/fredjn/etos/cli-go/pkg/client/v0/events"
	"github.com/fredjn/etos/cli-go/pkg/command"
	"github.com/fredjn/etos/cli-go/pkg/config"
	"github.com/spf13/cobra"
)

var testrunCmd = command.NewCommand(
	"testrun",
	"Operate on ETOS testruns",
	`Operate on ETOS testruns.

This command provides operations for managing ETOS testruns, such as
starting new testruns, checking their status, and retrieving results.`,
)

func init() {
	// Add the testrun command to the root command
	RootCmd.AddCommand(testrunCmd.Command)

	// Add version subcommands
	testrunCmd.AddCommand(newTestRunV0Command())
	// TODO: Add v1alpha command when implemented
}

// newTestRunV0Command creates the v0 subcommand
func newTestRunV0Command() *command.Command {
	cmd := command.NewCommand(
		"v0",
		"Manage ETOSv0 testruns",
		"Manage testruns using the ETOS v0 API",
	)

	// Add subcommands
	cmd.AddCommand(newTestRunListCommand(client.VersionV0))
	cmd.AddCommand(newTestRunGetCommand(client.VersionV0))
	cmd.AddCommand(newTestRunStartCommand())
	cmd.AddCommand(newTestRunStopCommand())
	cmd.AddCommand(newTestRunResultsCommand())

	return cmd
}

// newTestRunListCommand creates the list subcommand
func newTestRunListCommand(version client.Version) *command.Command {
	cmd := command.NewCommand(
		"list",
		"List testruns",
		"List all available testruns",
	)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Initialize configuration and client
		cfg := config.NewConfig()
		if err := cfg.Load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		var client client.VersionedClient
		switch version {
		case client.VersionV0:
			client = v0.NewClient(cfg)
		default:
			fmt.Printf("Unsupported version: %s\n", version)
			os.Exit(1)
		}

		ctx := context.Background()

		// List testruns
		testruns, err := client.ListTestRuns(ctx)
		if err != nil {
			fmt.Printf("Failed to list testruns: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Println("Testruns:")
		for _, tr := range testruns {
			fmt.Printf("- %s (Status: %s)\n", tr.ID, tr.Status)
		}
	}

	return cmd
}

// newTestRunGetCommand creates the get subcommand
func newTestRunGetCommand(version client.Version) *command.Command {
	cmd := command.NewCommand(
		"get",
		"Get testrun details",
		"Get detailed information about a specific testrun",
	)

	var testrunID string
	cmd.Flags().StringVarP(&testrunID, "id", "i", "", "Testrun ID")
	cmd.MarkFlagRequired("id")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Initialize configuration and client
		cfg := config.NewConfig()
		if err := cfg.Load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		var client client.VersionedClient
		switch version {
		case client.VersionV0:
			client = v0.NewClient(cfg)
		default:
			fmt.Printf("Unsupported version: %s\n", version)
			os.Exit(1)
		}

		ctx := context.Background()

		// Get testrun details
		testrun, err := client.GetTestRun(ctx, testrunID)
		if err != nil {
			fmt.Printf("Failed to get testrun: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Printf("Testrun Details:\n")
		fmt.Printf("ID: %s\n", testrun.ID)
		fmt.Printf("Status: %s\n", testrun.Status)
		fmt.Printf("Created: %s\n", testrun.CreatedAt)
		fmt.Printf("Updated: %s\n", testrun.UpdatedAt)
	}

	return cmd
}

// newTestRunStartCommand creates the start subcommand
func newTestRunStartCommand() *command.Command {
	cmd := command.NewCommand(
		"start",
		"Start a new testrun",
		"Start a new testrun with the specified parameters",
	)

	var (
		testSuiteID            string
		identity               string
		parentActivityID       string
		workspace              string
		artifactDir            string
		reportDir              string
		timeout                int
		iutProvider            string
		executionSpaceProvider string
		logAreaProvider        string
		dataset                []string
	)
	cmd.Flags().StringVarP(&testSuiteID, "test-suite", "s", "", "Test suite ID")
	cmd.Flags().StringVarP(&identity, "identity", "i", "", "Artifact created identity purl or ID")
	cmd.Flags().StringVarP(&parentActivityID, "parent-activity", "p", "", "Activity for the TERCC to link to")
	cmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Workspace to do all the work in")
	cmd.Flags().StringVarP(&artifactDir, "artifact-dir", "a", "", "Where test artifacts should be stored")
	cmd.Flags().StringVarP(&reportDir, "report-dir", "r", "", "Where test reports should be stored")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 3600, "Timeout in seconds")
	cmd.Flags().StringVar(&iutProvider, "iut-provider", "", "Which IUT provider to use")
	cmd.Flags().StringVar(&executionSpaceProvider, "execution-space-provider", "", "Which execution space provider to use")
	cmd.Flags().StringVar(&logAreaProvider, "log-area-provider", "", "Which log area provider to use")
	cmd.Flags().StringSliceVar(&dataset, "dataset", []string{}, "Additional dataset information to the environment provider")

	cmd.MarkFlagRequired("test-suite")
	cmd.MarkFlagRequired("identity")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Initialize configuration and client
		cfg := config.NewConfig()
		if err := cfg.Load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		client := v0.NewClient(cfg)
		ctx := context.Background()

		// Parse dataset
		var parsedDataset []map[string]interface{}
		for _, d := range dataset {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(d), &data); err != nil {
				fmt.Printf("Failed to parse dataset: %v\n", err)
				os.Exit(1)
			}
			parsedDataset = append(parsedDataset, data)
		}

		// Start testrun
		params := client.TestRunParams{
			TestSuiteID: testSuiteID,
			Identity:    identity,
			Timeout:     time.Duration(timeout) * time.Second,
			ProviderConfig: client.ProviderConfig{
				IUTProvider:            iutProvider,
				ExecutionSpaceProvider: executionSpaceProvider,
				LogAreaProvider:        logAreaProvider,
			},
			Dataset: parsedDataset,
		}

		if parentActivityID != "" {
			params.ParentActivityID = parentActivityID
		}
		if workspace != "" {
			params.Workspace = workspace
		}
		if artifactDir != "" {
			params.ArtifactDir = artifactDir
		}
		if reportDir != "" {
			params.ReportDir = reportDir
		}

		testrun, err := client.StartTestRun(ctx, params)
		if err != nil {
			fmt.Printf("Failed to start testrun: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Printf("Testrun started successfully:\n")
		fmt.Printf("ID: %s\n", testrun.ID)
		fmt.Printf("Status: %s\n", testrun.Status)
	}

	return cmd
}

// newTestRunStopCommand creates the stop subcommand
func newTestRunStopCommand() *command.Command {
	cmd := command.NewCommand(
		"stop",
		"Stop a testrun",
		"Stop a running testrun",
	)

	var testrunID string
	cmd.Flags().StringVarP(&testrunID, "id", "i", "", "Testrun ID")
	cmd.MarkFlagRequired("id")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Initialize configuration and client
		cfg := config.NewConfig()
		if err := cfg.Load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		client := v0.NewClient(cfg)
		ctx := context.Background()

		// Stop testrun
		if err := client.StopTestRun(ctx, testrunID); err != nil {
			fmt.Printf("Failed to stop testrun: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Testrun %s stopped successfully\n", testrunID)
	}

	return cmd
}

// newTestRunResultsCommand creates the results subcommand
func newTestRunResultsCommand() *command.Command {
	cmd := command.NewCommand(
		"results",
		"Get testrun results",
		"Get detailed test results for a specific testrun",
	)

	var (
		testrunID    string
		repository   string
	)
	cmd.Flags().StringVarP(&testrunID, "id", "i", "", "Testrun ID")
	cmd.Flags().StringVarP(&repository, "repository", "r", "", "Event repository URL")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("repository")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Initialize configuration and client
		cfg := config.NewConfig()
		if err := cfg.Load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			os.Exit(1)
		}

		client := v0.NewClient(cfg)
		ctx := context.Background()

		// Get testrun details
		testrun, err := client.GetTestRun(ctx, testrunID)
		if err != nil {
			fmt.Printf("Failed to get testrun: %v\n", err)
			os.Exit(1)
		}

		// Create event repository
		repo := events.NewGraphQLRepository(repository)

		// Get test results
		success, message, err := client.GetTestResults(ctx, testrunID, repo)
		if err != nil {
			fmt.Printf("Failed to get test results: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Printf("Test Results:\n")
		fmt.Printf("Success: %v\n", success)
		fmt.Printf("Message: %s\n", message)
	}

	return cmd
} 