// Package cmd implements all CLI commands using Cobra.
package cmd

import (
	"fmt"

	"github.com/enolalab/linear-cli/internal/api"
	"github.com/enolalab/linear-cli/internal/config"
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	// version is set at build time via ldflags.
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Global flags.
	flagPretty bool
	flagAPIKey string
)

// rootCmd is the base command.
var rootCmd = &cobra.Command{
	Use:   "linear-cli",
	Short: "CLI for Linear — designed for AI Agents",
	Long: `linear-cli is a command-line interface for the Linear project management tool.

It is designed primarily for AI Agents, providing:
  - JSON output by default (machine-readable)
  - Structured error responses with meaningful exit codes
  - No interactive prompts
  - Cursor-based pagination support

Authentication:
  Set LINEAR_API_KEY environment variable, or use 'linear-cli auth login --token <key>'.

Exit Codes:
  0  Success
  1  General error
  2  Invalid arguments
  3  Authentication error
  4  Resource not found
  5  Rate limit exceeded
  6  Network error
  7  Permission denied`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// versionCmd prints version information.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]string{
			"version": version,
			"commit":  commit,
			"date":    date,
		}
		output.PrintSuccess(data, nil)
	},
}

// Execute runs the root command and returns the exit code.
func Execute() int {
	// Initialize config.
	if err := config.Init(); err != nil {
		fmt.Fprintf(rootCmd.ErrOrStderr(), "warning: failed to load config: %v\n", err)
	}

	rootCmd.AddCommand(versionCmd)

	// Register all subcommands.
	registerAuthCommands()
	registerConfigCommands()
	registerTeamCommands()
	registerUserCommands()
	registerIssueCommands()
	registerCommentCommands()
	registerLabelCommands()
	registerStatusCommands()
	registerProjectCommands()
	registerCycleCommands()
	registerDocCommands()
	registerAttachmentCommands()

	if err := rootCmd.Execute(); err != nil {
		if cErr, ok := err.(*cliErr.CLIError); ok {
			return output.PrintError(cErr)
		}
		return output.PrintRawError(cliErr.ExitGeneral, cliErr.CodeGeneral, err.Error())
	}

	return cliErr.ExitSuccess
}

func init() {
	cobra.OnInitialize(initGlobals)

	rootCmd.PersistentFlags().BoolVar(&flagPretty, "pretty", false, "Pretty-print JSON output")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "Linear API key (overrides LINEAR_API_KEY env var)")
}

// initGlobals applies global flag values.
func initGlobals() {
	output.SetPretty(flagPretty)
}

// getAPIClient creates an API client using the resolved API key.
// Priority: --api-key flag > LINEAR_API_KEY env > config file.
func getAPIClient() (*api.Client, error) {
	key := flagAPIKey
	if key == "" {
		key = config.GetAPIKey()
	}
	if key == "" {
		return nil, cliErr.NewAuthError("no API key configured. Set LINEAR_API_KEY env var or run 'linear-cli auth login --token <key>'")
	}
	return api.NewClient(key), nil
}
