package users

import (
	"net/http"

	"github.com/akxcix/nomadcore/pkg/handlers"
	"github.com/akxcix/nomadcore/pkg/services/users"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	Service *users.Service
}

func New(s *users.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract Group ID from URL using Chi router
	username := chi.URLParam(r, "username")
	if username == "" {
		handlers.RespondWithError(w, r, ErrBadRequest)
		return
	}

	userProfile, err := h.Service.GetUserProfile(username)
	if err != nil {
		if err == users.ErrNoRowFound {
			handlers.RespondWithError(w, r, ErrDataNotFund.Wrap(err))
			return
		}
		handlers.RespondWithError(w, r, ErrInternalServerError.Wrap(err))
		return
	}

	userProfileDTO := UserProfileDTO{
		Username:  username,
		CreatedAt: userProfile.CreatedAt,
		UserID:    userProfile.ID,
		Bio:       userProfile.Bio,
		Interests: userProfile.Interests,
		Metadata:  userProfile.Metadata,
	}
	handlers.RespondWithData(w, r, userProfileDTO)
}
