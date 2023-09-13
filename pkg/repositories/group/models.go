package group

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

type Group struct {
	BaseTable
	UserId      string  `db:"user_id"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
}

type GroupDate struct {
	BaseTable
	FromDate time.Time `db:"from_date"`
	ToDate   time.Time `db:"to_date"`
	GroupID  uuid.UUID `db:"group_id"`
}

type GroupUser struct {
	BaseTable
	GroupID uuid.UUID `db:"group_id"`
	UserID  uuid.UUID `db:"user_id"`
}
