package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/enolalab/linear-cli/internal/config"
	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long: `Manage Linear issues. Subcommands: list, get, create, update, search.

Identifiers use the human-readable format TEAM-NUMBER (e.g., ENG-123).
The CLI automatically resolves these to UUIDs via the API.

Response fields for issue objects:
  id            UUID of the issue
  identifier    Human-readable identifier (e.g., ENG-123)
  title         Issue title
  description   Issue description (markdown, only in 'get')
  priority      Priority number (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)
  priorityLabel Priority as string ("Urgent", "High", etc.)
  estimate      Point estimate (nullable)
  state         { id, name, type, color } — workflow state
  assignee      { id, name, email } — assigned user (nullable)
  team          { id, name, key } — owning team
  labels        { nodes: [{ id, name, color }] } — applied labels
  cycle         { id, name, number } — sprint (nullable)
  project       { id, name } — parent project (nullable)
  url           Web URL to the issue
  createdAt     ISO 8601 timestamp
  updatedAt     ISO 8601 timestamp`,
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues with optional filters",
	Long: `List issues, optionally filtered by team, assignee, status, priority, or label.

Filters:
  --team       Team key (e.g., ENG). Falls back to default_team from config.
  --assignee   Assignee user ID or "me" for current user.
  --status     Workflow state name (e.g., "In Progress", "Done").
  --priority   Priority level: 0=No priority, 1=Urgent, 2=High, 3=Medium, 4=Low.
  --label      Label name to filter by.

Output JSON:
  { "success": true, "data": [ { "id": "...", "identifier": "ENG-123", ... } ], "pagination": { ... } }`,
	Example: `  linear-cli issue list --team ENG
  linear-cli issue list --team ENG --status "In Progress"
  linear-cli issue list --assignee me --priority 1
  linear-cli issue list --limit 10 --cursor "abc123"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")
		team, _ := cmd.Flags().GetString("team")
		assignee, _ := cmd.Flags().GetString("assignee")
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetInt("priority")
		label, _ := cmd.Flags().GetString("label")

		// Fall back to default team.
		if team == "" {
			team = config.GetDefaultTeam()
		}

		// Build filter.
		filter := map[string]interface{}{}
		if team != "" {
			filter["team"] = map[string]interface{}{
				"key": map[string]interface{}{"eq": team},
			}
		}
		if assignee != "" {
			if assignee == "me" {
				filter["assignee"] = map[string]interface{}{
					"isMe": map[string]interface{}{"eq": true},
				}
			} else {
				filter["assignee"] = map[string]interface{}{
					"id": map[string]interface{}{"eq": assignee},
				}
			}
		}
		if status != "" {
			filter["state"] = map[string]interface{}{
				"name": map[string]interface{}{"eq": status},
			}
		}
		if cmd.Flags().Changed("priority") {
			filter["priority"] = map[string]interface{}{
				"eq": priority,
			}
		}
		if label != "" {
			filter["labels"] = map[string]interface{}{
				"name": map[string]interface{}{"eq": label},
			}
		}

		query := `query($first: Int, $after: String, $filter: IssueFilter) {
			issues(first: $first, after: $after, filter: $filter, orderBy: updatedAt) {
				nodes {
					id identifier title
					priority priorityLabel
					estimate
					state { id name type color }
					assignee { id name email }
					team { id name key }
					labels { nodes { id name color } }
					cycle { id name number }
					project { id name }
					createdAt updatedAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}
		if cursor != "" {
			vars["after"] = cursor
		}
		if len(filter) > 0 {
			vars["filter"] = filter
		}

		var result struct {
			Issues struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"issues"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Issues.Nodes, &output.Pagination{
			HasNextPage: result.Issues.PageInfo.HasNextPage,
			EndCursor:   result.Issues.PageInfo.EndCursor,
		})
		return nil
	},
}

var issueGetCmd = &cobra.Command{
	Use:   "get <identifier>",
	Short: "Get issue by identifier (e.g., ENG-123)",
	Long: `Get detailed information about an issue by its human-readable identifier.

The identifier format is TEAM-NUMBER (e.g., ENG-123).
The response includes the full description and comments, which are not included in 'issue list'.

Response includes all standard issue fields plus:
  description   Full markdown description
  comments      { nodes: [{ id, body, createdAt, user: { id, name } }] }
  url           Direct link to issue in Linear web app

Output JSON:
  { "success": true, "data": { "id": "...", "identifier": "ENG-123", "title": "...", ... } }`,
	Example: `  linear-cli issue get ENG-123
  linear-cli issue get PRD-42 --pretty`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		identifier := args[0]
		teamKey, number, parseErr := parseIdentifier(identifier)
		if parseErr != nil {
			return parseErr
		}

		query := `query($filter: IssueFilter) {
			issues(filter: $filter) {
				nodes {
					id identifier title description
					priority priorityLabel
					estimate
					state { id name type color }
					assignee { id name email }
					team { id name key }
					labels { nodes { id name color } }
					cycle { id name number }
					project { id name }
					comments { nodes { id body createdAt user { id name } } }
					createdAt updatedAt
					url
				}
			}
		}`

		vars := map[string]interface{}{
			"filter": map[string]interface{}{
				"team": map[string]interface{}{
					"key": map[string]interface{}{"eq": teamKey},
				},
				"number": map[string]interface{}{
					"eq": number,
				},
			},
		}

		var result struct {
			Issues struct {
				Nodes []interface{} `json:"nodes"`
			} `json:"issues"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if len(result.Issues.Nodes) == 0 {
			return cliErr.NewNotFoundError("issue not found: " + identifier)
		}

		output.PrintSuccess(result.Issues.Nodes[0], nil)
		return nil
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long: `Create a new issue in Linear.

Required: --title and --team (or default_team in config).

The description can be provided via --description flag or read from a file via --description-file.

How to find IDs for flags:
  --assignee   Get user IDs from 'linear-cli user list'
  --label      Use label name directly (e.g., "bug") — resolved automatically
  --cycle      Get cycle IDs from 'linear-cli cycle list --team ENG'
  --project    Get project IDs from 'linear-cli project list'
  --state      Get state IDs from 'linear-cli status list --team ENG'
  --priority   Use integer directly: 0=None, 1=Urgent, 2=High, 3=Medium, 4=Low

Output JSON:
  { "success": true, "data": { "id": "...", "identifier": "ENG-124", "title": "...", "url": "..." } }`,
	Example: `  linear-cli issue create --title "Fix login bug" --team ENG --priority 1
  linear-cli issue create --title "Refactor auth" --team ENG --description "Details here"
  linear-cli issue create --title "Big task" --team ENG --description-file ./description.md
  linear-cli issue create --title "Sprint task" --team ENG --cycle CYCLE_ID --assignee USER_ID`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		team, _ := cmd.Flags().GetString("team")
		description, _ := cmd.Flags().GetString("description")
		descFile, _ := cmd.Flags().GetString("description-file")
		priority, _ := cmd.Flags().GetInt("priority")
		assignee, _ := cmd.Flags().GetString("assignee")
		labelName, _ := cmd.Flags().GetString("label")
		estimate, _ := cmd.Flags().GetFloat64("estimate")
		cycleID, _ := cmd.Flags().GetString("cycle")
		projectID, _ := cmd.Flags().GetString("project")
		stateID, _ := cmd.Flags().GetString("state")

		if title == "" {
			return cliErr.NewUsageError("--title is required")
		}

		if team == "" {
			team = config.GetDefaultTeam()
		}
		if team == "" {
			return cliErr.NewUsageError("--team is required (or set default_team in config)")
		}

		// Read description from file if specified.
		if descFile != "" {
			data, err := os.ReadFile(descFile)
			if err != nil {
				return cliErr.Wrap(cliErr.ExitUsage, cliErr.CodeUsage,
					fmt.Sprintf("failed to read description file: %s", descFile), err)
			}
			description = string(data)
		}

		// First, resolve team key to team ID.
		teamID, err := resolveTeamID(client, team)
		if err != nil {
			return err
		}

		input := map[string]interface{}{
			"title":  title,
			"teamId": teamID,
		}
		if description != "" {
			input["description"] = description
		}
		if cmd.Flags().Changed("priority") {
			input["priority"] = priority
		}
		if assignee != "" {
			input["assigneeId"] = assignee
		}
		if cmd.Flags().Changed("estimate") {
			input["estimate"] = estimate
		}
		if cycleID != "" {
			input["cycleId"] = cycleID
		}
		if projectID != "" {
			input["projectId"] = projectID
		}
		if stateID != "" {
			input["stateId"] = stateID
		}
		if labelName != "" {
			// Resolve label name to ID.
			labelID, err := resolveLabelID(client, teamID, labelName)
			if err != nil {
				return err
			}
			input["labelIds"] = []string{labelID}
		}

		query := `mutation($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				success
				issue {
					id identifier title
					priority priorityLabel
					state { id name }
					assignee { id name }
					team { id key }
					url
					createdAt
				}
			}
		}`

		vars := map[string]interface{}{"input": input}

		var result struct {
			IssueCreate struct {
				Success bool        `json:"success"`
				Issue   interface{} `json:"issue"`
			} `json:"issueCreate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.IssueCreate.Issue, nil)
		return nil
	},
}

var issueUpdateCmd = &cobra.Command{
	Use:   "update <identifier>",
	Short: "Update an existing issue",
	Long: `Update an issue by its identifier (e.g., ENG-123).

All fields are optional — only specified fields will be updated.

Status vs State:
  --status   Use workflow state NAME (e.g., "In Progress", "Done") — resolved automatically.
             Run 'linear-cli status list --team ENG' to see available state names.
  --state    Use workflow state UUID directly (skips resolution).

How to find IDs for flags:
  --assignee   Get user IDs from 'linear-cli user list'
  --label      Use label name directly (e.g., "bug") — resolved automatically
  --cycle      Get cycle IDs from 'linear-cli cycle list --team ENG'
  --project    Get project IDs from 'linear-cli project list'

Output JSON:
  { "success": true, "data": { "id": "...", "identifier": "ENG-123", ... } }`,
	Example: `  linear-cli issue update ENG-123 --status "In Progress"
  linear-cli issue update ENG-123 --priority 1 --assignee user-id
  linear-cli issue update ENG-123 --title "Updated title"
  linear-cli issue update ENG-123 --status Done --label bug`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		identifier := args[0]

		// Resolve identifier to issue ID.
		issueID, teamID, err := resolveIssueID(client, identifier)
		if err != nil {
			return err
		}

		input := map[string]interface{}{}

		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			input["title"] = v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input["description"] = v
		}
		if cmd.Flags().Changed("description-file") {
			f, _ := cmd.Flags().GetString("description-file")
			data, err := os.ReadFile(f)
			if err != nil {
				return cliErr.Wrap(cliErr.ExitUsage, cliErr.CodeUsage, "failed to read description file", err)
			}
			input["description"] = string(data)
		}
		if cmd.Flags().Changed("priority") {
			v, _ := cmd.Flags().GetInt("priority")
			input["priority"] = v
		}
		if cmd.Flags().Changed("assignee") {
			v, _ := cmd.Flags().GetString("assignee")
			input["assigneeId"] = v
		}
		if cmd.Flags().Changed("estimate") {
			v, _ := cmd.Flags().GetFloat64("estimate")
			input["estimate"] = v
		}
		if cmd.Flags().Changed("cycle") {
			v, _ := cmd.Flags().GetString("cycle")
			input["cycleId"] = v
		}
		if cmd.Flags().Changed("project") {
			v, _ := cmd.Flags().GetString("project")
			input["projectId"] = v
		}
		if cmd.Flags().Changed("label") {
			labelName, _ := cmd.Flags().GetString("label")
			labelID, err := resolveLabelID(client, teamID, labelName)
			if err != nil {
				return err
			}
			input["labelIds"] = []string{labelID}
		}
		if cmd.Flags().Changed("status") {
			statusName, _ := cmd.Flags().GetString("status")
			stateID, err := resolveStateID(client, teamID, statusName)
			if err != nil {
				return err
			}
			input["stateId"] = stateID
		}
		if cmd.Flags().Changed("state") {
			v, _ := cmd.Flags().GetString("state")
			input["stateId"] = v
		}

		if len(input) == 0 {
			return cliErr.NewUsageError("at least one field must be specified to update")
		}

		query := `mutation($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				success
				issue {
					id identifier title
					priority priorityLabel
					state { id name }
					assignee { id name }
					team { id key }
					url
					updatedAt
				}
			}
		}`

		vars := map[string]interface{}{
			"id":    issueID,
			"input": input,
		}

		var result struct {
			IssueUpdate struct {
				Success bool        `json:"success"`
				Issue   interface{} `json:"issue"`
			} `json:"issueUpdate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.IssueUpdate.Issue, nil)
		return nil
	},
}

var issueSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Full-text search issues",
	Long: `Search issues using Linear's full-text search.

The search query matches against issue title, description, and comments.
Results are ranked by relevance. Use 'issue list' with filters for structured queries.

Response includes: id, identifier, title, priority, priorityLabel, state, assignee, team, url, updatedAt.

Output JSON:
  { "success": true, "data": [ { "id": "...", "identifier": "ENG-123", "title": "..." } ], "pagination": { ... } }`,
	Example: `  linear-cli issue search "login bug"
  linear-cli issue search "authentication" --limit 5
  linear-cli issue search "NullPointerException stack trace"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		searchQuery := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		query := `query($term: String!, $first: Int) {
			searchIssues(term: $term, first: $first) {
				nodes {
					id identifier title
					priority priorityLabel
					state { id name type }
					assignee { id name }
					team { id key }
					url
					updatedAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{
			"term":  searchQuery,
			"first": limit,
		}

		var result struct {
			SearchIssues struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"searchIssues"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.SearchIssues.Nodes, &output.Pagination{
			HasNextPage: result.SearchIssues.PageInfo.HasNextPage,
			EndCursor:   result.SearchIssues.PageInfo.EndCursor,
		})
		return nil
	},
}

// parseIdentifier parses "ENG-123" into team key "ENG" and number 123.
func parseIdentifier(identifier string) (string, int, error) {
	parts := strings.SplitN(identifier, "-", 2)
	if len(parts) != 2 {
		return "", 0, cliErr.NewUsageError(
			fmt.Sprintf("invalid identifier format: %q (expected TEAM-NUMBER, e.g., ENG-123)", identifier))
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, cliErr.NewUsageError(
			fmt.Sprintf("invalid issue number in identifier: %q", identifier))
	}

	return strings.ToUpper(parts[0]), num, nil
}

// resolveTeamID resolves a team key (e.g., "ENG") to its UUID.
func resolveTeamID(client interface{ Execute(string, map[string]interface{}, interface{}) error }, teamKey string) (string, error) {
	query := `query($filter: TeamFilter) {
		teams(filter: $filter) {
			nodes { id }
		}
	}`

	vars := map[string]interface{}{
		"filter": map[string]interface{}{
			"key": map[string]interface{}{"eq": teamKey},
		},
	}

	var result struct {
		Teams struct {
			Nodes []struct {
				ID string `json:"id"`
			} `json:"nodes"`
		} `json:"teams"`
	}

	if err := client.Execute(query, vars, &result); err != nil {
		return "", err
	}

	if len(result.Teams.Nodes) == 0 {
		return "", cliErr.NewNotFoundError("team not found: " + teamKey)
	}

	return result.Teams.Nodes[0].ID, nil
}

// resolveIssueID resolves a human-readable identifier to issue UUID and team UUID.
func resolveIssueID(client interface{ Execute(string, map[string]interface{}, interface{}) error }, identifier string) (string, string, error) {
	teamKey, number, err := parseIdentifier(identifier)
	if err != nil {
		return "", "", err
	}

	query := `query($filter: IssueFilter) {
		issues(filter: $filter) {
			nodes { id team { id } }
		}
	}`

	vars := map[string]interface{}{
		"filter": map[string]interface{}{
			"team": map[string]interface{}{
				"key": map[string]interface{}{"eq": teamKey},
			},
			"number": map[string]interface{}{
				"eq": number,
			},
		},
	}

	var result struct {
		Issues struct {
			Nodes []struct {
				ID   string `json:"id"`
				Team struct {
					ID string `json:"id"`
				} `json:"team"`
			} `json:"nodes"`
		} `json:"issues"`
	}

	if err := client.Execute(query, vars, &result); err != nil {
		return "", "", err
	}

	if len(result.Issues.Nodes) == 0 {
		return "", "", cliErr.NewNotFoundError("issue not found: " + identifier)
	}

	return result.Issues.Nodes[0].ID, result.Issues.Nodes[0].Team.ID, nil
}

// resolveLabelID resolves a label name to its UUID within a team.
func resolveLabelID(client interface{ Execute(string, map[string]interface{}, interface{}) error }, teamID, labelName string) (string, error) {
	query := `query($filter: IssueLabelFilter) {
		issueLabels(filter: $filter) {
			nodes { id name }
		}
	}`

	vars := map[string]interface{}{
		"filter": map[string]interface{}{
			"name": map[string]interface{}{"eq": labelName},
			"team": map[string]interface{}{
				"id": map[string]interface{}{"eq": teamID},
			},
		},
	}

	var result struct {
		IssueLabels struct {
			Nodes []struct {
				ID string `json:"id"`
			} `json:"nodes"`
		} `json:"issueLabels"`
	}

	if err := client.Execute(query, vars, &result); err != nil {
		return "", err
	}

	if len(result.IssueLabels.Nodes) == 0 {
		return "", cliErr.NewNotFoundError("label not found: " + labelName)
	}

	return result.IssueLabels.Nodes[0].ID, nil
}

// resolveStateID resolves a workflow state name to its UUID within a team.
func resolveStateID(client interface{ Execute(string, map[string]interface{}, interface{}) error }, teamID, stateName string) (string, error) {
	query := `query($filter: WorkflowStateFilter) {
		workflowStates(filter: $filter) {
			nodes { id name }
		}
	}`

	vars := map[string]interface{}{
		"filter": map[string]interface{}{
			"name": map[string]interface{}{"eq": stateName},
			"team": map[string]interface{}{
				"id": map[string]interface{}{"eq": teamID},
			},
		},
	}

	var result struct {
		WorkflowStates struct {
			Nodes []struct {
				ID string `json:"id"`
			} `json:"nodes"`
		} `json:"workflowStates"`
	}

	if err := client.Execute(query, vars, &result); err != nil {
		return "", err
	}

	if len(result.WorkflowStates.Nodes) == 0 {
		return "", cliErr.NewNotFoundError("workflow state not found: " + stateName)
	}

	return result.WorkflowStates.Nodes[0].ID, nil
}

// cliErrNotFound is a helper for consistent not-found errors.
func cliErrNotFound(resource, identifier string) error {
	return cliErr.NewNotFoundError(fmt.Sprintf("%s not found: %s", resource, identifier))
}

func registerIssueCommands() {
	// issue list flags
	issueListCmd.Flags().String("team", "", "Filter by team key (e.g., ENG)")
	issueListCmd.Flags().String("assignee", "", "Filter by assignee ID or 'me'")
	issueListCmd.Flags().String("status", "", "Filter by workflow state name")
	issueListCmd.Flags().Int("priority", 0, "Filter by priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	issueListCmd.Flags().String("label", "", "Filter by label name")
	issueListCmd.Flags().Int("limit", 50, "Maximum number of results")
	issueListCmd.Flags().String("cursor", "", "Pagination cursor")

	// issue create flags
	issueCreateCmd.Flags().String("title", "", "Issue title (required)")
	issueCreateCmd.Flags().String("team", "", "Team key (e.g., ENG)")
	issueCreateCmd.Flags().String("description", "", "Issue description (markdown)")
	issueCreateCmd.Flags().String("description-file", "", "Read description from file")
	issueCreateCmd.Flags().Int("priority", 0, "Priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	issueCreateCmd.Flags().String("assignee", "", "Assignee user ID")
	issueCreateCmd.Flags().String("label", "", "Label name")
	issueCreateCmd.Flags().Float64("estimate", 0, "Point estimate")
	issueCreateCmd.Flags().String("cycle", "", "Cycle ID")
	issueCreateCmd.Flags().String("project", "", "Project ID")
	issueCreateCmd.Flags().String("state", "", "Workflow state ID")

	// issue update flags
	issueUpdateCmd.Flags().String("title", "", "New title")
	issueUpdateCmd.Flags().String("description", "", "New description")
	issueUpdateCmd.Flags().String("description-file", "", "Read new description from file")
	issueUpdateCmd.Flags().Int("priority", 0, "New priority")
	issueUpdateCmd.Flags().String("assignee", "", "New assignee user ID")
	issueUpdateCmd.Flags().String("label", "", "Label name to set")
	issueUpdateCmd.Flags().Float64("estimate", 0, "New estimate")
	issueUpdateCmd.Flags().String("cycle", "", "New cycle ID")
	issueUpdateCmd.Flags().String("project", "", "New project ID")
	issueUpdateCmd.Flags().String("status", "", "New status name (resolves to state ID)")
	issueUpdateCmd.Flags().String("state", "", "New workflow state ID (direct)")

	// issue search flags
	issueSearchCmd.Flags().Int("limit", 20, "Maximum number of results")

	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueGetCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueUpdateCmd)
	issueCmd.AddCommand(issueSearchCmd)
	rootCmd.AddCommand(issueCmd)
}
