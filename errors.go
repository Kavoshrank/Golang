package kavoshrank

import (
	"fmt"
)

// APIError represents an error returned by the KavoshRank API.
type APIError struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("kavoshrank api error (status %d): %s, details: %v", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("kavoshrank api error (status %d): %s", e.StatusCode, e.Message)
}

// AuthenticationError represents a 401 or 403 error.
type AuthenticationError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("kavoshrank authentication error (status %d): %s", e.StatusCode, e.Message)
}

// ValidationError represents a 400 validation error.
type ValidationError struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors,omitempty"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("kavoshrank validation error (status %d): %s, details: %v", e.StatusCode, e.Message, e.Errors)
}

// RateLimitError represents a 429 rate limit error.
type RateLimitError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after,omitempty"`
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("kavoshrank rate limit exceeded (status %d): %s, retry after %d seconds", e.StatusCode, e.Message, e.RetryAfter)
	}
	return fmt.Sprintf("kavoshrank rate limit exceeded (status %d): %s", e.StatusCode, e.Message)
}
