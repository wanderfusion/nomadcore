package passport

import "github.com/google/uuid"

type GetUsersFromIDsResponse struct {
	Status int       `json:"status"`
	Data   []UserDTO `json:"data"`
}

type UserDTO struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	ProfilePic string    `json:"profilePic"`
}
