// Package api provides a GraphQL HTTP client for the Linear API.
// It uses raw HTTP POST requests instead of a typed GraphQL client
// for simplicity and flexibility.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	cliErr "github.com/enolalab/linear-cli/internal/errors"
)

const (
	defaultBaseURL = "https://api.linear.app/graphql"
	defaultTimeout = 30 * time.Second
	userAgent      = "linear-cli"
)

// Version is set at build time via ldflags.
var Version = "dev"

// Client is the Linear GraphQL API client.
type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// graphQLRequest is the JSON body sent to the GraphQL endpoint.
type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// graphQLResponse is the raw response from the GraphQL endpoint.
type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphQLError  `json:"errors,omitempty"`
}

// graphQLError represents a single GraphQL error.
type graphQLError struct {
	Message    string                   `json:"message"`
	Extensions *graphQLErrorExtensions  `json:"extensions,omitempty"`
}

// graphQLErrorExtensions may contain error codes or rate limit info.
type graphQLErrorExtensions struct {
	Code     string `json:"code,omitempty"`
	Type     string `json:"type,omitempty"`
}

// NewClient creates a new Linear API client.
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
	}
}

// Execute sends a GraphQL query and unmarshals the result into the provided target.
// The target should be a pointer to a struct matching the "data" field of the GraphQL response.
func (c *Client) Execute(query string, variables map[string]interface{}, target interface{}) error {
	reqBody := graphQLRequest{
		Query:     query,
		Variables: variables,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return cliErr.NewGeneralError("failed to marshal request", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return cliErr.NewGeneralError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", userAgent, Version))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return cliErr.NewNetworkError("failed to connect to Linear API", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return cliErr.NewNetworkError("failed to read response", err)
	}

	// Handle HTTP-level errors.
	if resp.StatusCode == http.StatusUnauthorized {
		return cliErr.NewAuthError("invalid or expired API key")
	}
	if resp.StatusCode == http.StatusForbidden {
		return cliErr.NewPermissionError("insufficient permissions for this operation")
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return cliErr.NewRateLimitError("API rate limit exceeded, try again later")
	}
	if resp.StatusCode >= 500 {
		return cliErr.NewNetworkError(
			fmt.Sprintf("Linear API returned HTTP %d", resp.StatusCode), nil,
		)
	}

	// Parse GraphQL response.
	var gqlResp graphQLResponse
	if err := json.Unmarshal(respBytes, &gqlResp); err != nil {
		return cliErr.NewGeneralError("failed to parse API response", err)
	}

	// Handle GraphQL-level errors.
	if len(gqlResp.Errors) > 0 {
		return c.handleGraphQLErrors(gqlResp.Errors)
	}

	// Unmarshal data into target.
	if target != nil && gqlResp.Data != nil {
		if err := json.Unmarshal(gqlResp.Data, target); err != nil {
			return cliErr.NewGeneralError("failed to parse API data", err)
		}
	}

	return nil
}

// handleGraphQLErrors converts GraphQL errors to typed CLI errors.
func (c *Client) handleGraphQLErrors(errors []graphQLError) *cliErr.CLIError {
	if len(errors) == 0 {
		return nil
	}

	// Collect all error messages.
	messages := make([]string, len(errors))
	for i, e := range errors {
		messages[i] = e.Message
	}
	combined := strings.Join(messages, "; ")

	// Check first error for specific types.
	first := errors[0]
	if first.Extensions != nil {
		switch first.Extensions.Code {
		case "AUTHENTICATION_ERROR":
			return cliErr.NewAuthError(combined)
		case "FORBIDDEN":
			return cliErr.NewPermissionError(combined)
		case "RATELIMITED":
			return cliErr.NewRateLimitError(combined)
		}
	}

	// Check message content for common patterns.
	lowerMsg := strings.ToLower(combined)
	if strings.Contains(lowerMsg, "not found") || strings.Contains(lowerMsg, "does not exist") {
		return cliErr.NewNotFoundError(combined)
	}

	return cliErr.NewGeneralError(combined, nil)
}
