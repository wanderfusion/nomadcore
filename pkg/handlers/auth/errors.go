package auth

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/errors"
)

var (
	ErrInvalidRefreshToken = errors.NewAPIError("invalid refresh token", http.StatusBadRequest)
	ErrInvalidAuthToken    = errors.NewAPIError("invalid auth token", http.StatusBadRequest)
)
