package cmd

import (
	"github.com/enolalab/linear-cli/internal/config"
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a configuration value. Available keys:
  - default_team: Default team key (e.g., ENG) used when --team is not specified.

Output JSON:
  { "success": true, "data": { "key": "...", "value": "..." } }`,
	Example: `  linear-cli config set default_team ENG`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, value := args[0], args[1]

		if err := config.Set(key, value); err != nil {
			return cliErr.NewGeneralError("failed to save config", err)
		}

		output.PrintSuccess(map[string]string{
			"key":   key,
			"value": value,
		}, nil)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Long: `Get a configuration value by key.

Output JSON:
  { "success": true, "data": { "key": "...", "value": "..." } }`,
	Example: `  linear-cli config get default_team`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := config.Get(key)

		output.PrintSuccess(map[string]string{
			"key":   key,
			"value": value,
		}, nil)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	Long: `List all configuration values (API key is redacted).

Output JSON:
  { "success": true, "data": { "default_team": "ENG", "api_key": "***" } }`,
	Example: `  linear-cli config list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all := config.GetAll()

		// Redact API key for security.
		if _, ok := all[config.KeyAPIKey]; ok {
			val, isStr := all[config.KeyAPIKey].(string)
			if isStr && len(val) > 8 {
				all[config.KeyAPIKey] = val[:8] + "***"
			} else if isStr && len(val) > 0 {
				all[config.KeyAPIKey] = "***"
			}
		}

		output.PrintSuccess(all, nil)
		return nil
	},
}

func registerConfigCommands() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
