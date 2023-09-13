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

type AddUsersToGroupReq struct {
	GroupID uuid.UUID   `json:"groupID"`
	UserIDs []uuid.UUID `json:"userIDs"`
}

// Res --------------------------------------------------------------------------------------------

type GetGroupsRes struct {
	Groups []GroupDTO `json:"groups"`
}

type GetGroupDetailsRes struct {
	Group     GroupDTO  `json:"group"`
	GroupDate []DateDTO `json:"dates"`
	GroupUser []UserDTO `json:"users"`
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

type UserDTO struct {
	ID uuid.UUID `json:"id"`
}
