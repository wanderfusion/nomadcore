package errors

import (
	"errors"
	"fmt"
)

type Error interface {
	GetMsg() string
	GetErr() error
	GetStatusCode() int
	Wrap(error) Error
}

type APIError struct {
	Msg    string
	Err    error
	Status int
}

func (e APIError) GetMsg() string {
	return e.Msg
}

func (e APIError) GetErr() error {
	return e.Err
}

func (e APIError) GetStatusCode() int {
	return e.Status
}

func (e APIError) Wrap(err error) Error {
	return APIError{
		Msg:    e.Msg,
		Err:    fmt.Errorf("%w: %v", err, e.Err),
		Status: e.Status,
	}
}

func NewAPIError(msg string, status int) APIError {
	return APIError{
		Msg:    msg,
		Err:    errors.New(msg),
		Status: status,
	}
}
