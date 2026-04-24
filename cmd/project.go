package cmd

import (
	"os"

	"github.com/enolalab/linear-cli/internal/config"
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List projects",
	Long: `List projects, optionally filtered by team.

Output JSON:
  { "success": true, "data": [ { "id": "...", "name": "...", "state": "..." } ], "pagination": { ... } }`,
	Example: `  linear-cli project list
  linear-cli project list --team ENG`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		team, _ := cmd.Flags().GetString("team")

		query := `query($first: Int, $after: String, $filter: ProjectFilter) {
			projects(first: $first, after: $after, filter: $filter) {
				nodes {
					id name description
					state slugId url icon color
					progress
					startDate targetDate
					lead { id name }
					teams { nodes { id key } }
					createdAt updatedAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}
		if cursor != "" {
			vars["after"] = cursor
		}
		if team != "" {
			teamID, resolveErr := resolveTeamID(client, team)
			if resolveErr != nil {
				return resolveErr
			}
			vars["filter"] = map[string]interface{}{
				"accessibleTeams": map[string]interface{}{
					"id": map[string]interface{}{"eq": teamID},
				},
			}
		}

		var result struct {
			Projects struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"projects"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Projects.Nodes, &output.Pagination{
			HasNextPage: result.Projects.PageInfo.HasNextPage,
			EndCursor:   result.Projects.PageInfo.EndCursor,
		})
		return nil
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get <project-id>",
	Short: "Get project by ID",
	Long: `Get detailed information about a project.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "progress": 0.75, ... } }`,
	Example: `  linear-cli project get abc-123`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			project(id: $id) {
				id name description
				state slugId url icon color
				progress
				startDate targetDate
				lead { id name email }
				members { nodes { id name email } }
				teams { nodes { id name key } }
				issues { nodes { id identifier title state { name } } }
				createdAt updatedAt
			}
		}`

		vars := map[string]interface{}{"id": args[0]}

		var result struct {
			Project interface{} `json:"project"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if result.Project == nil {
			return cliErr.NewNotFoundError("project not found: " + args[0])
		}

		output.PrintSuccess(result.Project, nil)
		return nil
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long: `Create a new project in Linear.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", "url": "..." } }`,
	Example: `  linear-cli project create --name "Q2 Roadmap" --team ENG
  linear-cli project create --name "Auth Revamp" --team ENG --description-file ./plan.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		descFile, _ := cmd.Flags().GetString("description-file")
		team, _ := cmd.Flags().GetString("team")

		if name == "" {
			return cliErr.NewUsageError("--name is required")
		}

		if team == "" {
			team = config.GetDefaultTeam()
		}
		if team == "" {
			return cliErr.NewUsageError("--team is required")
		}

		if descFile != "" {
			data, readErr := os.ReadFile(descFile)
			if readErr != nil {
				return cliErr.Wrap(cliErr.ExitUsage, cliErr.CodeUsage, "failed to read description file", readErr)
			}
			description = string(data)
		}

		teamID, resolveErr := resolveTeamID(client, team)
		if resolveErr != nil {
			return resolveErr
		}

		input := map[string]interface{}{
			"name":    name,
			"teamIds": []string{teamID},
		}
		if description != "" {
			input["description"] = description
		}

		query := `mutation($input: ProjectCreateInput!) {
			projectCreate(input: $input) {
				success
				project { id name state url createdAt }
			}
		}`

		vars := map[string]interface{}{"input": input}

		var result struct {
			ProjectCreate struct {
				Success bool        `json:"success"`
				Project interface{} `json:"project"`
			} `json:"projectCreate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.ProjectCreate.Project, nil)
		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update <project-id>",
	Short: "Update an existing project",
	Long: `Update a project by its ID.

Output JSON:
  { "success": true, "data": { "id": "...", "name": "...", ... } }`,
	Example: `  linear-cli project update abc-123 --name "Updated Name"
  linear-cli project update abc-123 --state completed`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		input := map[string]interface{}{}

		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input["name"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input["description"] = v
		}
		if cmd.Flags().Changed("state") {
			v, _ := cmd.Flags().GetString("state")
			input["state"] = v
		}

		if len(input) == 0 {
			return cliErr.NewUsageError("at least one field must be specified to update")
		}

		query := `mutation($id: String!, $input: ProjectUpdateInput!) {
			projectUpdate(id: $id, input: $input) {
				success
				project { id name state url updatedAt }
			}
		}`

		vars := map[string]interface{}{
			"id":    args[0],
			"input": input,
		}

		var result struct {
			ProjectUpdate struct {
				Success bool        `json:"success"`
				Project interface{} `json:"project"`
			} `json:"projectUpdate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.ProjectUpdate.Project, nil)
		return nil
	},
}

func registerProjectCommands() {
	projectListCmd.Flags().String("team", "", "Filter by team key")
	projectListCmd.Flags().Int("limit", 50, "Maximum number of results")
	projectListCmd.Flags().String("cursor", "", "Pagination cursor")

	projectCreateCmd.Flags().String("name", "", "Project name (required)")
	projectCreateCmd.Flags().String("team", "", "Team key")
	projectCreateCmd.Flags().String("description", "", "Project description")
	projectCreateCmd.Flags().String("description-file", "", "Read description from file")

	projectUpdateCmd.Flags().String("name", "", "New name")
	projectUpdateCmd.Flags().String("description", "", "New description")
	projectUpdateCmd.Flags().String("state", "", "New state (planned, started, paused, completed, cancelled)")

	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	rootCmd.AddCommand(projectCmd)
}
