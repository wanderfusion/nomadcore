package users

import (
	"errors"

	"github.com/akxcix/nomadcore/pkg/services"
)

var (
	// 5xx ----------------------------------------------------------------------------------------
	ErrFailedDBWrite    services.ServiceError = errors.New("something went wrong while writing to the database")
	ErrFailedDBRead     services.ServiceError = errors.New("something went wrong while reading from the database")
	ErrFailedClientCall services.ServiceError = errors.New("something went wrong while calling a client")

	// 4xx ----------------------------------------------------------------------------------------
	ErrInvalidRequest services.ServiceError = errors.New("the request is invalid")
	ErrUserForbidden  services.ServiceError = errors.New("the user is forbidden to access this resource")
	ErrNoRowFound     services.ServiceError = errors.New("no row found")
)
