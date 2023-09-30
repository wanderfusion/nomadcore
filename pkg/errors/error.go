// Package errors defines custom error types and interfaces for API error handling.
package errors

import (
	"errors"
	"fmt"
)

// Error is an interface for API-specific errors.
type Error interface {
	GetMsg() string     // Returns the error message
	GetErr() error      // Returns the underlying error object
	GetStatusCode() int // Returns the HTTP status code associated with the error
	Wrap(error) Error   // Wraps an existing error with additional context
}

// APIError is a struct that implements the Error interface.
type APIError struct {
	Msg    string // Human-readable message describing the error
	Err    error  // Underlying error object
	Status int    // HTTP status code for the error
}

// GetMsg returns the error message.
func (e APIError) GetMsg() string {
	return e.Msg
}

// GetErr returns the underlying error object.
func (e APIError) GetErr() error {
	return e.Err
}

// GetStatusCode returns the HTTP status code.
func (e APIError) GetStatusCode() int {
	return e.Status
}

// Wrap wraps an existing error with the APIError.
func (e APIError) Wrap(err error) Error {
	return APIError{
		Msg:    e.Msg,
		Err:    fmt.Errorf("%w: %v", err, e.Err),
		Status: e.Status,
	}
}

// NewAPIError creates a new APIError.
// Parameters:
// - msg: The error message.
// - status: The HTTP status code for the error.
func NewAPIError(msg string, status int) APIError {
	return APIError{
		Msg:    msg,
		Err:    errors.New(msg),
		Status: status,
	}
}
