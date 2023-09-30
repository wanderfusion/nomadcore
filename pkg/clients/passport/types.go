// Package passport defines data structures for handling user data and responses.
package passport

import "github.com/google/uuid"

// GetUsersFromIDsResponse represents the structure of the API response for fetching multiple users.
type GetUsersFromIDsResponse struct {
	Status int       `json:"status"` // HTTP status code
	Data   []UserDTO `json:"data"`   // Array of user data
}

// UserDTO is a Data Transfer Object representing the structure of a user.
type UserDTO struct {
	ID         uuid.UUID `json:"id"`         // Unique identifier for the user
	Username   string    `json:"username"`   // Username of the user
	ProfilePic string    `json:"profilePic"` // URL to the profile picture of the user
}
