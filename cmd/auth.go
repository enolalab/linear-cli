package cmd

import (
	"github.com/enolalab/linear-cli/internal/config"
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save API key to config",
	Long: `Save a Linear API key to the local config file (~/.config/linear-cli/config.yaml).

Alternatively, set the LINEAR_API_KEY environment variable.

Output JSON:
  { "success": true, "data": { "message": "API key saved", "configDir": "..." } }`,
	Example: `  linear-cli auth login --token lin_api_xxxxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			return cliErr.NewUsageError("--token flag is required")
		}

		if err := config.Set(config.KeyAPIKey, token); err != nil {
			return cliErr.NewGeneralError("failed to save API key", err)
		}

		output.PrintSuccess(map[string]string{
			"message":   "API key saved successfully",
			"configDir": config.ConfigDir(),
		}, nil)
		return nil
	},
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authenticated user",
	Long: `Display information about the currently authenticated Linear user.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "email": "..." } }`,
	Example: `  linear-cli auth whoami
  LINEAR_API_KEY=lin_api_xxx linear-cli auth whoami`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query { viewer { id name email displayName active admin } }`

		var result struct {
			Viewer struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Email       string `json:"email"`
				DisplayName string `json:"displayName"`
				Active      bool   `json:"active"`
				Admin       bool   `json:"admin"`
			} `json:"viewer"`
		}

		if err := client.Execute(query, nil, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Viewer, nil)
		return nil
	},
}

func registerAuthCommands() {
	authLoginCmd.Flags().String("token", "", "Linear API key (required)")

	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authWhoamiCmd)
	rootCmd.AddCommand(authCmd)
}
