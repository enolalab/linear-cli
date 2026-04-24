// Package errors provides typed errors with meaningful exit codes.
// AI Agents use exit codes to programmatically handle errors
// without parsing error messages.
package errors

import "fmt"

// Exit codes — each maps to a specific error category.
const (
	ExitSuccess    = 0 // Operation completed successfully
	ExitGeneral    = 1 // Unspecified error
	ExitUsage      = 2 // Invalid arguments or usage error
	ExitAuth       = 3 // Authentication failure (missing/invalid API key)
	ExitNotFound   = 4 // Requested resource not found
	ExitRateLimit  = 5 // API rate limit exceeded
	ExitNetwork    = 6 // Network connectivity error
	ExitPermission = 7 // Insufficient permissions
)

// Error code strings used in JSON error responses.
const (
	CodeGeneral    = "GENERAL_ERROR"
	CodeUsage      = "USAGE_ERROR"
	CodeAuth       = "AUTH_ERROR"
	CodeNotFound   = "NOT_FOUND"
	CodeRateLimit  = "RATE_LIMITED"
	CodeNetwork    = "NETWORK_ERROR"
	CodePermission = "PERMISSION_DENIED"
)

// CLIError is a typed error that carries an exit code and error code string.
type CLIError struct {
	ExitCode int
	Code     string
	Message  string
	Err      error
}

func (e *CLIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *CLIError) Unwrap() error {
	return e.Err
}

// New creates a CLIError with the given exit code, code string, and message.
func New(exitCode int, code string, message string) *CLIError {
	return &CLIError{
		ExitCode: exitCode,
		Code:     code,
		Message:  message,
	}
}

// Wrap creates a CLIError wrapping an existing error.
func Wrap(exitCode int, code string, message string, err error) *CLIError {
	return &CLIError{
		ExitCode: exitCode,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}

// Convenience constructors for common error types.

func NewUsageError(message string) *CLIError {
	return New(ExitUsage, CodeUsage, message)
}

func NewAuthError(message string) *CLIError {
	return New(ExitAuth, CodeAuth, message)
}

func NewNotFoundError(message string) *CLIError {
	return New(ExitNotFound, CodeNotFound, message)
}

func NewRateLimitError(message string) *CLIError {
	return New(ExitRateLimit, CodeRateLimit, message)
}

func NewNetworkError(message string, err error) *CLIError {
	return Wrap(ExitNetwork, CodeNetwork, message, err)
}

func NewPermissionError(message string) *CLIError {
	return New(ExitPermission, CodePermission, message)
}

func NewGeneralError(message string, err error) *CLIError {
	return Wrap(ExitGeneral, CodeGeneral, message, err)
}
