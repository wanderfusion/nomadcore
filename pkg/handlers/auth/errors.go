package auth

import (
	"net/http"

	"github.com/wanderfusion/nomadcore/pkg/errors"
)

var (
	ErrInvalidRefreshToken = errors.NewAPIError("invalid refresh token", http.StatusBadRequest)
	ErrInvalidAuthToken    = errors.NewAPIError("invalid auth token", http.StatusBadRequest)
)
