package group

import (
	"time"

	"github.com/google/uuid"
)

// Req --------------------------------------------------------------------------------------------
type CreateGroupReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddDatesToGroupReq struct {
	GroupID uuid.UUID `json:"groupID"`
	Dates   DateDTO   `json:"dates"`
}

// Res --------------------------------------------------------------------------------------------

type GetGroupsRes struct {
	Groups []GroupDTO `json:"groups"`
}

// DTOs -------------------------------------------------------------------------------------------
type GroupDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Dates       []DateDTO `json:"dates"`
}

type DateDTO struct {
	ID   uuid.UUID `json:"id"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
