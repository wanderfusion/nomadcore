package calendar

import "github.com/google/uuid"

// Req --------------------------------------------------------------------------------------------
type CreateCalendarReq struct {
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

type JwtVerifyRequest struct {
	Jwt string `json:"jwt"`
}

// Res --------------------------------------------------------------------------------------------

type GetPublicCalendarsRes struct {
	Calendars []CalendarDTO `json:"calendars"`
}

// DTOs -------------------------------------------------------------------------------------------
type CalendarDTO struct {
	ID         uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	Visibility string    `json:"visibility"`
}
