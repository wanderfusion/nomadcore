package users

import "github.com/google/uuid"

func (db *Database) GetUserProfileByID(userID uuid.UUID) (UserProfile, error) {
	query := `
		SELECT
		*
		FROM 
		user_profiles
		WHERE
		user_id = $1	
	`

	userProfile := UserProfile{}
	err := db.db.Get(&userProfile, query, userID)
	if err != nil {
		return userProfile, err
	}

	return userProfile, nil
}
