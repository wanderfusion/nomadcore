package users

import (
	"time"

	"github.com/google/uuid"
)

// base types -------------------------------------------------------------------------------------
type BaseTable struct {
	ID        uuid.UUID `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserProfile struct {
	BaseTable
	UserID    uuid.UUID `db:"user_id"`
	Bio       string    `db:"bio"`
	Interests string    `db:"interests"`
	Metadata  string    `db:"metadata"`
}
