package groq

import (
	"fmt"
)

// APIError represents an error returned by the Groq API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Groq API error (status %d): %s", e.StatusCode, e.Message)
}