package cmd

import (
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace users",
	Long: `List all users in the workspace.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "...", "email": "..." } ], "pagination": { ... } }`,
	Example: `  linear-cli user list
  linear-cli user list --limit 20`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")

		query := `query($first: Int, $after: String) {
			users(first: $first, after: $after) {
				nodes {
					id name email displayName
					active admin guest
					avatarUrl
					createdAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}
		if cursor != "" {
			vars["after"] = cursor
		}

		var result struct {
			Users struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"users"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Users.Nodes, &output.Pagination{
			HasNextPage: result.Users.PageInfo.HasNextPage,
			EndCursor:   result.Users.PageInfo.EndCursor,
		})
		return nil
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get user by ID",
	Long: `Get user details by their Linear user ID.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "email": "..." } }`,
	Example: `  linear-cli user get abc123-def456`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			user(id: $id) {
				id name email displayName
				active admin guest
				avatarUrl
				createdAt
			}
		}`

		vars := map[string]interface{}{"id": args[0]}

		var result struct {
			User interface{} `json:"user"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if result.User == nil {
			return cliErr.NewNotFoundError("user not found: " + args[0])
		}

		output.PrintSuccess(result.User, nil)
		return nil
	},
}

var userMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current authenticated user (alias for 'auth whoami')",
	Long: `Display the currently authenticated user. Same as 'linear-cli auth whoami'.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "email": "..." } }`,
	Example: `  linear-cli user me`,
	RunE:    authWhoamiCmd.RunE,
}

func registerUserCommands() {
	userListCmd.Flags().Int("limit", 50, "Maximum number of results")
	userListCmd.Flags().String("cursor", "", "Pagination cursor")

	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userMeCmd)
	rootCmd.AddCommand(userCmd)
}
