package cmd

import (
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage teams",
}

var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all teams in the workspace",
	Long: `List all teams accessible to the authenticated user.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "...", "key": "..." } ], "pagination": { ... } }`,
	Example: `  linear-cli team list
  linear-cli team list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")

		query := `query($first: Int, $after: String) {
			teams(first: $first, after: $after) {
				nodes {
					id name key description
					private
					timezone
					issueCount
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}
		if cursor != "" {
			vars["after"] = cursor
		}

		var result struct {
			Teams struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"teams"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Teams.Nodes, &output.Pagination{
			HasNextPage: result.Teams.PageInfo.HasNextPage,
			EndCursor:   result.Teams.PageInfo.EndCursor,
		})
		return nil
	},
}

var teamGetCmd = &cobra.Command{
	Use:   "get <team-key>",
	Short: "Get team details by key (e.g., ENG)",
	Long: `Get detailed information about a team by its key.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "key": "...", "members": [...] } }`,
	Example: `  linear-cli team get ENG`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		teamKey := args[0]

		query := `query($filter: TeamFilter) {
			teams(filter: $filter) {
				nodes {
					id name key description
					private timezone issueCount
					members {
						nodes { id name email displayName active }
					}
					states {
						nodes { id name type position color }
					}
					labels {
						nodes { id name color }
					}
				}
			}
		}`

		vars := map[string]interface{}{
			"filter": map[string]interface{}{
				"key": map[string]interface{}{"eq": teamKey},
			},
		}

		var result struct {
			Teams struct {
				Nodes []interface{} `json:"nodes"`
			} `json:"teams"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if len(result.Teams.Nodes) == 0 {
			return cliErrNotFound("team", teamKey)
		}

		output.PrintSuccess(result.Teams.Nodes[0], nil)
		return nil
	},
}

func registerTeamCommands() {
	teamListCmd.Flags().Int("limit", 50, "Maximum number of results")
	teamListCmd.Flags().String("cursor", "", "Pagination cursor (endCursor from previous response)")

	teamCmd.AddCommand(teamListCmd)
	teamCmd.AddCommand(teamGetCmd)
	rootCmd.AddCommand(teamCmd)
}
