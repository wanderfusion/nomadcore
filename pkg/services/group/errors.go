package group

import (
	"errors"

	"github.com/akxcix/nomadcore/pkg/services"
)

var (
	// 5xx ----------------------------------------------------------------------------------------
	ErrFailedDBWrite services.ServiceError = errors.New("something went wrong while writing to the database")
	ErrFailedDBRead  services.ServiceError = errors.New("something went wrong while reading from the database")

	// 4xx ----------------------------------------------------------------------------------------
	ErrInvalidRequest services.ServiceError = errors.New("the request is invalid")
)
