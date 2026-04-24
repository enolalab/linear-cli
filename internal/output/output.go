// Package output provides a consistent JSON envelope for all CLI output.
// Every response follows the same schema so AI Agents can reliably parse output.
package output

import (
	"encoding/json"
	"fmt"
	"os"

	cliErr "github.com/enolalab/linear-cli/internal/errors"
)

// Response is the standard JSON envelope for all CLI output.
type Response struct {
	Success    bool         `json:"success"`
	Data       interface{}  `json:"data,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
}

// Pagination contains cursor-based pagination info.
type Pagination struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor,omitempty"`
}

// ErrorDetail contains structured error information.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// pretty controls whether JSON output is indented.
var pretty bool

// SetPretty enables or disables indented JSON output.
func SetPretty(v bool) {
	pretty = v
}

// marshal converts a value to JSON bytes, respecting the pretty setting.
func marshal(v interface{}) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// PrintSuccess outputs a successful response with data and optional pagination.
func PrintSuccess(data interface{}, pagination *Pagination) {
	resp := Response{
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
	b, err := marshal(resp)
	if err != nil {
		// Fallback: if we can't marshal the response, print a raw error.
		fmt.Fprintf(os.Stderr, "fatal: failed to marshal response: %v\n", err)
		os.Exit(cliErr.ExitGeneral)
	}
	fmt.Println(string(b))
}

// PrintError outputs an error response and returns the appropriate exit code.
func PrintError(cErr *cliErr.CLIError) int {
	resp := Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    cErr.Code,
			Message: cErr.Error(),
		},
	}
	b, err := marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: failed to marshal error response: %v\n", err)
		return cliErr.ExitGeneral
	}
	fmt.Println(string(b))
	return cErr.ExitCode
}

// PrintRawError creates a CLIError from raw parameters and prints it.
// Returns the exit code.
func PrintRawError(exitCode int, code string, message string) int {
	return PrintError(cliErr.New(exitCode, code, message))
}
