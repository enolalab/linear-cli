package cmd

import (
	"os"

	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage issue comments",
}

var commentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments on an issue",
	Long: `List all comments on an issue by its identifier.

Output JSON:
  { "success": true, "data": [ { "id": "...", "body": "...", "user": {...}, "createdAt": "..." } ] }`,
	Example: `  linear-cli comment list --issue ENG-123`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		issueIdentifier, _ := cmd.Flags().GetString("issue")
		if issueIdentifier == "" {
			return cliErr.NewUsageError("--issue is required")
		}

		issueID, _, err := resolveIssueID(client, issueIdentifier)
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			issue(id: $id) {
				comments {
					nodes {
						id body
						user { id name email }
						createdAt updatedAt
					}
				}
			}
		}`

		vars := map[string]interface{}{"id": issueID}

		var result struct {
			Issue struct {
				Comments struct {
					Nodes []interface{} `json:"nodes"`
				} `json:"comments"`
			} `json:"issue"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Issue.Comments.Nodes, nil)
		return nil
	},
}

var commentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a comment on an issue",
	Long: `Add a comment to an issue. The body can be provided via --body flag or --body-file.

Output JSON:
  { "success": true, "data": { "id": "...", "body": "..." } }`,
	Example: `  linear-cli comment create --issue ENG-123 --body "This is fixed in PR #456"
  linear-cli comment create --issue ENG-123 --body-file ./review-notes.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		issueIdentifier, _ := cmd.Flags().GetString("issue")
		body, _ := cmd.Flags().GetString("body")
		bodyFile, _ := cmd.Flags().GetString("body-file")

		if issueIdentifier == "" {
			return cliErr.NewUsageError("--issue is required")
		}

		if bodyFile != "" {
			data, readErr := os.ReadFile(bodyFile)
			if readErr != nil {
				return cliErr.Wrap(cliErr.ExitUsage, cliErr.CodeUsage, "failed to read body file", readErr)
			}
			body = string(data)
		}

		if body == "" {
			return cliErr.NewUsageError("--body or --body-file is required")
		}

		issueID, _, err := resolveIssueID(client, issueIdentifier)
		if err != nil {
			return err
		}

		query := `mutation($input: CommentCreateInput!) {
			commentCreate(input: $input) {
				success
				comment {
					id body
					user { id name }
					createdAt
				}
			}
		}`

		vars := map[string]interface{}{
			"input": map[string]interface{}{
				"issueId": issueID,
				"body":    body,
			},
		}

		var result struct {
			CommentCreate struct {
				Success bool        `json:"success"`
				Comment interface{} `json:"comment"`
			} `json:"commentCreate"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.CommentCreate.Comment, nil)
		return nil
	},
}

func registerCommentCommands() {
	commentListCmd.Flags().String("issue", "", "Issue identifier (e.g., ENG-123) (required)")
	commentCreateCmd.Flags().String("issue", "", "Issue identifier (e.g., ENG-123) (required)")
	commentCreateCmd.Flags().String("body", "", "Comment body (markdown)")
	commentCreateCmd.Flags().String("body-file", "", "Read comment body from file")

	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentCreateCmd)
	rootCmd.AddCommand(commentCmd)
}
