package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cliErr "github.com/enolalab/linear-cli/internal/errors"
	"github.com/enolalab/linear-cli/internal/output"
	"github.com/spf13/cobra"
)

var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage file attachments",
}

var attachmentUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file and attach it to an issue",
	Long: `Upload a file to Linear's storage and create an attachment on an issue.

This is a multi-step process:
  1. Request a presigned upload URL from Linear (fileUpload mutation)
  2. Upload the file to the presigned URL (HTTP PUT)
  3. Create an attachment linking the uploaded file to the issue (attachmentCreate mutation)

Output JSON:
  { "success": true, "data": { "id": "...", "url": "...", "title": "...", "issueIdentifier": "..." } }`,
	Example: `  linear-cli attachment upload --issue ENG-123 --file ./screenshot.png
  linear-cli attachment upload --issue ENG-123 --file ./report.pdf --title "Bug Report"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAPIClient()
		if err != nil {
			return err
		}

		issueIdentifier, _ := cmd.Flags().GetString("issue")
		filePath, _ := cmd.Flags().GetString("file")
		title, _ := cmd.Flags().GetString("title")

		if issueIdentifier == "" {
			return cliErr.NewUsageError("--issue is required")
		}
		if filePath == "" {
			return cliErr.NewUsageError("--file is required")
		}

		// Read the file.
		fileData, readErr := os.ReadFile(filePath)
		if readErr != nil {
			return cliErr.Wrap(cliErr.ExitUsage, cliErr.CodeUsage,
				fmt.Sprintf("failed to read file: %s", filePath), readErr)
		}

		fileName := filepath.Base(filePath)
		fileSize := len(fileData)
		contentType := detectContentType(filePath, fileData)

		if title == "" {
			title = fileName
		}

		// Resolve issue identifier to ID.
		issueID, _, resolveErr := resolveIssueID(client, issueIdentifier)
		if resolveErr != nil {
			return resolveErr
		}

		// Step 1: Request upload URL.
		uploadQuery := `mutation($contentType: String!, $filename: String!, $size: Int!) {
			fileUpload(contentType: $contentType, filename: $filename, size: $size) {
				uploadFile {
					uploadUrl
					assetUrl
					contentType
					filename
					size
				}
				success
			}
		}`

		uploadVars := map[string]interface{}{
			"contentType": contentType,
			"filename":    fileName,
			"size":        fileSize,
		}

		var uploadResult struct {
			FileUpload struct {
				Success    bool `json:"success"`
				UploadFile struct {
					UploadUrl   string `json:"uploadUrl"`
					AssetUrl    string `json:"assetUrl"`
					ContentType string `json:"contentType"`
					Filename    string `json:"filename"`
					Size        int    `json:"size"`
				} `json:"uploadFile"`
			} `json:"fileUpload"`
		}

		if err := client.Execute(uploadQuery, uploadVars, &uploadResult); err != nil {
			return err
		}

		if !uploadResult.FileUpload.Success {
			return cliErr.NewGeneralError("failed to get upload URL from Linear", nil)
		}

		uploadURL := uploadResult.FileUpload.UploadFile.UploadUrl
		assetURL := uploadResult.FileUpload.UploadFile.AssetUrl

		// Step 2: Upload file to presigned URL.
		httpClient := &http.Client{Timeout: 60 * time.Second}
		putReq, reqErr := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileData))
		if reqErr != nil {
			return cliErr.NewGeneralError("failed to create upload request", reqErr)
		}
		putReq.Header.Set("Content-Type", contentType)
		putReq.ContentLength = int64(fileSize)

		putResp, putErr := httpClient.Do(putReq)
		if putErr != nil {
			return cliErr.NewNetworkError("failed to upload file", putErr)
		}
		defer func() { _ = putResp.Body.Close() }()

		if putResp.StatusCode >= 400 {
			body, _ := io.ReadAll(putResp.Body)
			return cliErr.NewGeneralError(
				fmt.Sprintf("file upload failed with HTTP %d: %s", putResp.StatusCode, string(body)), nil)
		}

		// Step 3: Create attachment on the issue.
		attachQuery := `mutation($input: AttachmentCreateInput!) {
			attachmentCreate(input: $input) {
				success
				attachment {
					id title url
					creator { id name }
					createdAt
				}
			}
		}`

		attachVars := map[string]interface{}{
			"input": map[string]interface{}{
				"issueId": issueID,
				"title":   title,
				"url":     assetURL,
			},
		}

		var attachResult struct {
			AttachmentCreate struct {
				Success    bool        `json:"success"`
				Attachment interface{} `json:"attachment"`
			} `json:"attachmentCreate"`
		}

		if err := client.Execute(attachQuery, attachVars, &attachResult); err != nil {
			return err
		}

		// Combine upload + attachment info.
		responseData := map[string]interface{}{
			"attachment":      attachResult.AttachmentCreate.Attachment,
			"assetUrl":        assetURL,
			"fileName":        fileName,
			"fileSize":        fileSize,
			"contentType":     contentType,
			"issueIdentifier": issueIdentifier,
		}

		output.PrintSuccess(responseData, nil)
		return nil
	},
}

// detectContentType detects MIME type from file extension and content.
func detectContentType(filePath string, data []byte) string {
	// Try extension first.
	ext := filepath.Ext(filePath)
	if ext != "" {
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			return mimeType
		}
	}

	// Fall back to content sniffing.
	return http.DetectContentType(data)
}

func registerAttachmentCommands() {
	attachmentUploadCmd.Flags().String("issue", "", "Issue identifier (e.g., ENG-123) (required)")
	attachmentUploadCmd.Flags().String("file", "", "Path to file to upload (required)")
	attachmentUploadCmd.Flags().String("title", "", "Attachment title (defaults to filename)")

	attachmentCmd.AddCommand(attachmentUploadCmd)
	rootCmd.AddCommand(attachmentCmd)
}
