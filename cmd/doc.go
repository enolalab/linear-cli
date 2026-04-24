package cmd

import (
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Manage documents",
}

var docListCmd = &cobra.Command{
	Use:   "list",
	Short: "List documents",
	Long: `List documents in the workspace.

Output JSON:
  { "success": true, "data": [ { "id": "...", "title": "...", "content": "..." } ], "pagination": { ... } }`,
	Example: `  linear-cli doc list
  linear-cli doc list --limit 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		limit, _ := cmd.Flags().GetInt("limit")
		cursor, _ := cmd.Flags().GetString("cursor")

		query := `query($first: Int, $after: String) {
			documents(first: $first, after: $after) {
				nodes {
					id title slugId
					icon color
					content
					project { id name }
					creator { id name }
					updatedBy { id name }
					createdAt updatedAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{"first": limit}
		if cursor != "" {
			vars["after"] = cursor
		}

		var result struct {
			Documents struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"documents"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.Documents.Nodes, &output.Pagination{
			HasNextPage: result.Documents.PageInfo.HasNextPage,
			EndCursor:   result.Documents.PageInfo.EndCursor,
		})
		return nil
	},
}

var docGetCmd = &cobra.Command{
	Use:   "get <doc-id>",
	Short: "Get document by ID",
	Long: `Get detailed information about a document, including its full content.

Output JSON:
  { "success": true, "data": { "id": "...", "title": "...", "content": "..." } }`,
	Example: `  linear-cli doc get abc-123`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		query := `query($id: String!) {
			document(id: $id) {
				id title slugId
				icon color
				content
				project { id name }
				creator { id name email }
				updatedBy { id name }
				createdAt updatedAt
			}
		}`

		vars := map[string]interface{}{"id": args[0]}

		var result struct {
			Document interface{} `json:"document"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		if result.Document == nil {
			return cliErrNotFound("document", args[0])
		}

		output.PrintSuccess(result.Document, nil)
		return nil
	},
}

var docSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search documents",
	Long: `Search documents using Linear's search.

Output JSON:
  { "success": true, "data": [ { "id": "...", "title": "...", ... } ], "pagination": { ... } }`,
	Example: `  linear-cli doc search "onboarding"
  linear-cli doc search "architecture" --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		searchQuery := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		query := `query($term: String!, $first: Int) {
			searchDocuments(term: $term, first: $first) {
				nodes {
					id title slugId
					icon color
					content
					project { id name }
					creator { id name }
					createdAt updatedAt
				}
				pageInfo { hasNextPage endCursor }
			}
		}`

		vars := map[string]interface{}{
			"term":  searchQuery,
			"first": limit,
		}

		var result struct {
			SearchDocuments struct {
				Nodes    []interface{} `json:"nodes"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"searchDocuments"`
		}

		if err := client.Execute(query, vars, &result); err != nil {
			return err
		}

		output.PrintSuccess(result.SearchDocuments.Nodes, &output.Pagination{
			HasNextPage: result.SearchDocuments.PageInfo.HasNextPage,
			EndCursor:   result.SearchDocuments.PageInfo.EndCursor,
		})
		return nil
	},
}

func registerDocCommands() {
	docListCmd.Flags().Int("limit", 20, "Maximum number of results")
	docListCmd.Flags().String("cursor", "", "Pagination cursor")

	docSearchCmd.Flags().Int("limit", 20, "Maximum number of results")

	docCmd.AddCommand(docListCmd)
	docCmd.AddCommand(docGetCmd)
	docCmd.AddCommand(docSearchCmd)
	rootCmd.AddCommand(docCmd)
}
