package users

import (
	"net/http"

	"github.com/wanderfusion/nomadcore/pkg/errors"
)

var (
	ErrInvalidContext      = errors.NewAPIError("somthing went wrong", http.StatusInternalServerError)
	ErrBadRequest          = errors.NewAPIError("bad request", http.StatusBadRequest)
	ErrInternalServerError = errors.NewAPIError("something went wrong", http.StatusInternalServerError)
	ErrDataNotFund         = errors.NewAPIError("requested data not found", http.StatusNotFound)
)
