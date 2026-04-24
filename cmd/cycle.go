package cmd

import (
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Manage cycles (sprints)",
}

var cycleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cycles for a team",
	Long: `List all cycles (sprints) for a team.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "...", "number": 1, ... } ], "pagination": { ... } }`,
	Example: `  linear-cli cycle list --team ENG`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		team, _ := cmd.Flags().GetString("team")
		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")

		if team == "" {
			return cliErr.NewUsageError("--team is required")
		}

		teamID, resolveErr := resolveTeamID(client, team)
		if resolveErr != nil {
			return resolveErr
		}

		query := `query($first: Int, $after: String, $filter: CycleFilter) {
			cycles(first: $first, after: $after, filter: $filter) {
				nodes {
					id name number description
					startsAt endsAt
					completedAt
					progress
					scopeCompleted scope
					issueCountHistory completedIssueCountHistory
					team { id key }
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{
			"first": limit,
			"filter": map[string]interface{}{
				"team": map[string]interface{}{
					"id": map[string]interface{}{"eq": teamID},
				},
			},
		}
		if cursor != "" {
			vars["after"] = cursor
		}

		var result struct {
			Cycles struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"cycles"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Cycles.Nodes, &output.Pagination{
			HasNextPage: result.Cycles.PageInfo.HasNextPage,
			EndCursor:   result.Cycles.PageInfo.EndCursor,
		})
		return nil
	},
}

var cycleGetCmd = &cobra.Command{
	Use:   "get <cycle-id>",
	Short: "Get cycle by ID (includes issues)",
	Long: `Get detailed information about a cycle, including its issues.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "issues": [...] } }`,
	Example: `  linear-cli cycle get abc-123`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			cycle(id: $id) {
				id name number description
				startsAt endsAt completedAt
				progress
				scopeCompleted scope
				team { id key }
				issues {
					nodes {
						id identifier title
						priority priorityLabel
						state { name type }
						assignee { id name }
						estimate
					}
				}
			}
		}`

		vars := map[string]interface{}{"id": args[0]}

		var result struct {
			Cycle interface{} `json:"cycle"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if result.Cycle == nil {
			return cliErr.NewNotFoundError("cycle not found: " + args[0])
		}

		output.PrintSuccess(result.Cycle, nil)
		return nil
	},
}

func registerCycleCommands() {
	cycleListCmd.Flags().String("team", "", "Team key (required)")
	cycleListCmd.Flags().Int("limit", 20, "Maximum number of results")
	cycleListCmd.Flags().String("cursor", "", "Pagination cursor")

	cycleCmd.AddCommand(cycleListCmd)
	cycleCmd.AddCommand(cycleGetCmd)
	rootCmd.AddCommand(cycleCmd)
}
