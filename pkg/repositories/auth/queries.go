package auth

import (
	"github.com/google/uuid"
)

func (db *Database) CreateCalendar(userId uuid.UUID, name, visibility string) error {
	query := `
		INSERT INTO public.calendars (user_id, name, visibility) VALUES ($1, $2, $3)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, userId, name, visibility)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) GetPublicCalendars(userId uuid.UUID) ([]Calendars, error) {
	calendars := []Calendars{}
	query := `
		SELECT * FROM public.calendars WHERE user_id = $1 and visibility = 'public'
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.Select(&calendars, query, userId)
	if err != nil {
		return nil, err
	}

	return calendars, nil
}

// func (db *Database) UpdateUserProfile(user User) error {
// 	query := `
//         UPDATE public.users
//         SET
//             username = CASE WHEN $2::text != '' THEN $2::text ELSE username END,
//             profile_picture = CASE WHEN $3::text != '' THEN $3::text ELSE profile_picture END,
//             updated_at = NOW()
//         WHERE id = $1::uuid
//     `

// 	tx := db.db.MustBegin()
// 	_, err := tx.Exec(query, user.ID, user.Username, user.ProfilePic)
// 	if err != nil {
// 		return err
// 	}
// 	return tx.Commit()
// }
