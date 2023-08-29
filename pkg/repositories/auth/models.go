package auth

import (
	"time"

	"github.com/google/uuid"
)

type BaseTable struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Calendars struct {
	BaseTable
	UserId     string `db:"user_id"`
	Name       string `db:"name"`
	Visibility string `db:"visibility"`
}

type Dates struct {
	BaseTable
	FromDate   time.Time `db:"from_date"`
	ToDate     time.Time `db:"from_date"`
	CalendarId uuid.UUID `db:"calendar_id"`
}
