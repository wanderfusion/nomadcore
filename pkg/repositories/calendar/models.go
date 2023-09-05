package calendar

import (
	"time"

	"github.com/google/uuid"
)

// base types -------------------------------------------------------------------------------------
type Visibility string

const (
	Public  Visibility = "public"
	Private Visibility = "private"
)

type BaseTable struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Calendar struct {
	BaseTable
	UserId     string     `db:"user_id"`
	Name       string     `db:"name"`
	Visibility Visibility `db:"visibility"`
}

type Date struct {
	BaseTable
	FromDate   time.Time `db:"from_date"`
	ToDate     time.Time `db:"to_date"`
	CalendarId uuid.UUID `db:"calendar_id"`
}
