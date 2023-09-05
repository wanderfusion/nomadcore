package group

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (db *Database) CreateGroup(userId uuid.UUID, name, description string) error {
	query := `
		INSERT 
			INTO 
				public.groups (user_id, name, description) 
			VALUES 
				($1, $2, $3)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, userId, name, description)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) GetGroups(userID uuid.UUID) ([]Group, error) {
	groups := []Group{}
	query := `
		SELECT 
			* 
		FROM 
			public.groups 
		WHERE 
			user_id = $1 
	`

	// QueryRow still works, but now we're scanning into multiple variables
	err := db.db.Select(&groups, query, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (db *Database) GetDates(calendarIDs []uuid.UUID) ([]Date, error) {
	dates := []Date{}
	query := `
		SELECT 
			*
		FROM
			public.dates
		WHERE
			group_id in (?)
		LIMIT
			50
	`

	q, vs, err := sqlx.In(query, calendarIDs)
	if err != nil {
		return nil, err
	}

	q = db.db.Rebind(q)

	err = db.db.Select(&dates, q, vs...)
	if err != nil {
		return nil, err
	}

	return dates, nil
}

func (db *Database) AddDatesToGroup(userID, groupID uuid.UUID, from, to time.Time) error {
	query := `
	INSERT INTO 
		dates (from_date, to_date, group_id)
	SELECT 
		$1, $2, $3
	WHERE 
		EXISTS (
			SELECT 
				1
			FROM 
				groups
			WHERE 
				id = $3 
				AND user_id = $4
		);
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, from, to, groupID, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
