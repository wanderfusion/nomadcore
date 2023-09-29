package users

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/errors"
)

var (
	ErrInvalidContext      = errors.NewAPIError("somthing went wrong", http.StatusInternalServerError)
	ErrBadRequest          = errors.NewAPIError("bad request", http.StatusBadRequest)
	ErrInternalServerError = errors.NewAPIError("something went wrong", http.StatusInternalServerError)
)
