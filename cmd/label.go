package cmd

import (
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage issue labels",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issue labels",
	Long: `List issue labels, optionally filtered by team.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "Bug", "color": "#ff0000" } ] }`,
	Example: `  linear-cli label list
  linear-cli label list --team ENG`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		team, _ := cmd.Flags().GetString("team")
		limit, _ := cmd.Flags().GetInt("limit")

		query := `query($first: Int, $filter: IssueLabelFilter) {
			issueLabels(first: $first, filter: $filter) {
				nodes { id name color description }
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}

		if team != "" {
			teamID, resolveErr := resolveTeamID(client, team)
			if resolveErr != nil {
				return resolveErr
			}
			vars["filter"] = map[string]interface{}{
				"team": map[string]interface{}{
					"id": map[string]interface{}{"eq": teamID},
				},
			}
		}

		var result struct {
			IssueLabels struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"issueLabels"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.IssueLabels.Nodes, &output.Pagination{
			HasNextPage: result.IssueLabels.PageInfo.HasNextPage,
			EndCursor:   result.IssueLabels.PageInfo.EndCursor,
		})
		return nil
	},
}

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Long: `Create a new issue label.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "color": "..." } }`,
	Example: `  linear-cli label create --name "Bug" --color "#ff0000"
  linear-cli label create --name "Feature" --team ENG --color "#00ff00"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")
		team, _ := cmd.Flags().GetString("team")

		if name == "" {
			return cliErr.NewUsageError("--name is required")
		}

		input := map[string]interface{}{"name": name}
		if color != "" {
			input["color"] = color
		}
		if team != "" {
			teamID, resolveErr := resolveTeamID(client, team)
			if resolveErr != nil {
				return resolveErr
			}
			input["teamId"] = teamID
		}

		query := `mutation($input: IssueLabelCreateInput!) {
			issueLabelCreate(input: $input) {
				success
				issueLabel { id name color description }
			}
		}`

		vars := map[string]interface{}{"input": input}

		var result struct {
			IssueLabelCreate struct {
				Success    bool        `json:"success"`
				IssueLabel interface{} `json:"issueLabel"`
			} `json:"issueLabelCreate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.IssueLabelCreate.IssueLabel, nil)
		return nil
	},
}

func registerLabelCommands() {
	labelListCmd.Flags().String("team", "", "Filter by team key")
	labelListCmd.Flags().Int("limit", 50, "Maximum number of results")

	labelCreateCmd.Flags().String("name", "", "Label name (required)")
	labelCreateCmd.Flags().String("color", "", "Label color (hex, e.g., #ff0000)")
	labelCreateCmd.Flags().String("team", "", "Team key to scope label to")

	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelCreateCmd)
	rootCmd.AddCommand(labelCmd)
}
