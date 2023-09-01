package calendar

import (
	"context"
	"net/http"
	"strings"

	"github.com/akxcix/nomadcore/pkg/handlers"
	"github.com/akxcix/nomadcore/pkg/services/calendar"
	"github.com/google/uuid"
)

type Handlers struct {
	Service *calendar.Service
}

func New(s *calendar.Service) *Handlers {
	h := Handlers{
		Service: s,
	}

	return &h
}

func (h *Handlers) CreateCalendar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusBadRequest)
		return
	}
	var req CreateCalendarReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	err := h.Service.CalRepo.CreateCalendar(userID, req.Name, req.Visibility)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, "update successful")
}

func (h *Handlers) GetPublicCalendars(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusBadRequest)
		return
	}

	calendars, err := h.Service.CalRepo.GetPublicCalendars(userID)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	calendarDtos := make([]CalendarDTO, 0)
	for _, cal := range calendars {
		calDto := CalendarDTO{
			ID:         cal.ID,
			Name:       cal.Name,
			Visibility: cal.Visibility,
		}

		calendarDtos = append(calendarDtos, calDto)
	}

	res := GetPublicCalendarsRes{
		Calendars: calendarDtos,
	}

	handlers.RespondWithData(w, r, res)
}

func (h *Handlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		// Check if token exists
		if bearerToken == "" {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}

		// Extract token from header
		tokenParts := strings.Split(bearerToken, " ")
		if len(tokenParts) != 2 {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]

		// Validate token
		claims, isValid := h.Service.ValidateJwt(tokenString)
		if claims == nil || !isValid {
			handlers.RespondWithError(w, r, ErrInvalidJwt, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.UserIdContextKey, claims.ID)

		// If token is valid, forward to the actual handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
