package unstructured

import (
	"errors"
	"fmt"
)

// HTTPValidationError represents the structure of validation error responses
// returned by the API when a 422 status code is returned.
type HTTPValidationError struct {
	Detail []*ValidationError `json:"detail"`
}

func (e *HTTPValidationError) Error() string {
	errs := make([]error, len(e.Detail))
	for i, err := range e.Detail {
		errs[i] = err
	}

	return fmt.Sprintf("%d validation errors: %s", len(e.Detail), errors.Join(errs...))
}

// ValidationError represents a single validation error within the HTTPValidationError response.
type ValidationError struct {
	// Location is an array that can contain strings or integers indicating
	// where the validation error occurred (e.g., field names, array indices).
	Location []any `json:"loc"`
	// Message is a string describing the validation error.
	Message string `json:"msg"`
	// Type is a string indicating the type of error.
	Type string `json:"type"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s at %v: %s", e.Type, e.Location, e.Message)
}
