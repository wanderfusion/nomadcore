package calendar

import (
	"time"

	"github.com/google/uuid"
)

// Req --------------------------------------------------------------------------------------------
type CreateCalendarReq struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type AddDatesToCalendarReq struct {
	CalendarID uuid.UUID `json:"calendarId"`
	Dates      DateDTO   `json:"dates"`
}

// Res --------------------------------------------------------------------------------------------

type GetPublicCalendarsRes struct {
	Calendars []CalendarDTO `json:"calendars"`
}

// DTOs -------------------------------------------------------------------------------------------
type CalendarDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Visibility string    `json:"visibility"`
}

type DateDTO struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
