package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/akxcix/nomadcore/pkg/handlers"
	"github.com/akxcix/nomadcore/pkg/services/auth"
	"github.com/google/uuid"
)

type Handlers struct {
	Service *auth.Service
}

func New(s *auth.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			handlers.RespondWithError(w, r, ErrInvalidAuthToken.Wrap(err))
			return
		}

		claims, isValid := h.Service.ValidateJwt(token)
		if !isValid {
			handlers.RespondWithError(w, r, ErrInvalidAuthToken.Wrap(err))
			return
		}

		ctx := enrichContextWithUserID(r.Context(), claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractToken(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("missing token")
	}

	tokenParts := strings.Split(bearerToken, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("invalid token format")
	}

	return tokenParts[1], nil
}

func enrichContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, handlers.UserIdContextKey, userID)
}
