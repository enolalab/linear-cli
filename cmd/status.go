package cmd

import (
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Manage workflow states",
}

var statusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflow states for a team",
	Long: `List all workflow states (e.g., Backlog, Todo, In Progress, Done, Canceled) for a team.

State types: backlog, unstarted, started, completed, cancelled.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "In Progress", "type": "started", "color": "#..." } ] }`,
	Example: `  linear-cli status list --team ENG`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		team, _ := cmd.Flags().GetString("team")
		if team == "" {
			return cliErr.NewUsageError("--team is required")
		}

		teamID, resolveErr := resolveTeamID(client, team)
		if resolveErr != nil {
			return resolveErr
		}

		query := `query($filter: WorkflowStateFilter) {
			workflowStates(filter: $filter) {
				nodes {
					id name type position color description
					team { id key }
				}
			}
		}`

		vars := map[string]interface{}{
			"filter": map[string]interface{}{
				"team": map[string]interface{}{
					"id": map[string]interface{}{"eq": teamID},
				},
			},
		}

		var result struct {
			WorkflowStates struct {
				Nodes []interface{} `json:"nodes"`
			} `json:"workflowStates"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.WorkflowStates.Nodes, nil)
		return nil
	},
}

var statusGetCmd = &cobra.Command{
	Use:   "get <state-id>",
	Short: "Get workflow state by ID",
	Long: `Get a specific workflow state by its ID.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "type": "..." } }`,
	Example: `  linear-cli status get abc-123-def`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			workflowState(id: $id) {
				id name type position color description
				team { id key }
			}
		}`

		vars := map[string]interface{}{"id": args[0]}

		var result struct {
			WorkflowState interface{} `json:"workflowState"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if result.WorkflowState == nil {
			return cliErr.NewNotFoundError("workflow state not found: " + args[0])
		}

		output.PrintSuccess(result.WorkflowState, nil)
		return nil
	},
}

func registerStatusCommands() {
	statusListCmd.Flags().String("team", "", "Team key (required)")

	statusCmd.AddCommand(statusListCmd)
	statusCmd.AddCommand(statusGetCmd)
	rootCmd.AddCommand(statusCmd)
}
