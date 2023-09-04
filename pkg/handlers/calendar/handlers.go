package calendar

import (
	"net/http"

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
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}
	var req CreateCalendarReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	msg, err := h.Service.CreateCalendar(userID, req.Name, req.Visibility)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}

func (h *Handlers) GetPublicCalendars(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}

	calendars, err := h.Service.GetCalendars(userID, "public")
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	calendarDtos := make([]CalendarDTO, 0)
	for _, cal := range calendars {
		calDto := CalendarDTO{
			ID:         cal.ID,
			Name:       cal.Name,
			Visibility: string(cal.Visibility),
		}

		calendarDtos = append(calendarDtos, calDto)
	}

	res := GetPublicCalendarsRes{
		Calendars: calendarDtos,
	}

	handlers.RespondWithData(w, r, res)
}

func (h *Handlers) AddDatesToCalendar(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(handlers.UserIdContextKey).(uuid.UUID)
	if !ok {
		handlers.RespondWithError(w, r, ErrContextInvalid, http.StatusInternalServerError)
		return
	}
	var req AddDatesToCalendarReq
	if err := handlers.FromRequest(r, &req); err != nil {
		handlers.RespondWithError(w, r, err, http.StatusBadRequest)
		return
	}

	dates := calendar.Dates{
		From: req.Dates.From,
		To:   req.Dates.To,
	}
	msg, err := h.Service.AddDatesToCalendar(userID, req.CalendarID, dates)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	handlers.RespondWithData(w, r, msg)
}
