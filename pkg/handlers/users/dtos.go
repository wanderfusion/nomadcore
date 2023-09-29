package users

import (
	"time"

	"github.com/google/uuid"
)

type UserProfileDTO struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UserID    uuid.UUID `json:"userId"`
	Bio       string    `json:"bio"`
	Interests string    `json:"interests"`
	Metadata  string    `json:"metadata"`
}
