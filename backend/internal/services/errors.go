package services

import "fmt"

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

// ConflictError represents a resource conflict error
type ConflictError struct {
	Resource string
	Message  string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict error: %s - %s", e.Resource, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Resource: message,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Resource: message,
	}
}
