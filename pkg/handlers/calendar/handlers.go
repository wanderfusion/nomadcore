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

	calendarIds := make([]uuid.UUID, 0)
	for _, calendar := range calendars {
		calendarIds = append(calendarIds, calendar.ID)
	}

	dates, err := h.Service.GetDates(calendarIds)
	if err != nil {
		handlers.RespondWithError(w, r, err, http.StatusInternalServerError)
		return
	}

	dateMap := make(map[uuid.UUID][]DateDTO)
	for _, date := range dates {
		if _, exists := dateMap[date.CalendarId]; !exists {
			dateMap[date.CalendarId] = make([]DateDTO, 0)
		}

		dateSlice := dateMap[date.CalendarId]
		dateSlice = append(dateSlice, DateDTO{
			ID:   date.ID,
			From: date.FromDate,
			To:   date.ToDate,
		})
		dateMap[date.CalendarId] = dateSlice
	}

	calendarDTOs := make([]CalendarDTO, 0)
	for _, cal := range calendars {
		dateDtos := dateMap[cal.ID]

		calDto := CalendarDTO{
			ID:         cal.ID,
			Name:       cal.Name,
			Visibility: string(cal.Visibility),
			Dates:      dateDtos,
		}

		calendarDTOs = append(calendarDTOs, calDto)
	}

	res := GetPublicCalendarsRes{
		Calendars: calendarDTOs,
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
