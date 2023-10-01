package group

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (db *Database) CreateGroup(userID uuid.UUID, name, description string) (uuid.UUID, error) {
	query := `
		INSERT INTO public.groups (user_id, name, description) 
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var newGroupID uuid.UUID
	err := db.db.QueryRow(query, userID, name, description).Scan(&newGroupID)
	if err != nil {
		return uuid.Nil, err
	}
	return newGroupID, nil
}

func (db *Database) GetGroups(userID uuid.UUID) ([]Group, error) {
	groups := []Group{}
	query := `
		SELECT 
			g.* 
		FROM 
			public.groups g
		JOIN
			public.group_users gu ON g.id = gu.group_id
		WHERE 
			gu.user_id = $1 
	`

	err := db.db.Select(&groups, query, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (db *Database) GetGroupWithDetails(groupId uuid.UUID) (*Group, []GroupDate, []GroupUser, error) {
	var group Group
	var groupDates []GroupDate
	var groupUsers []GroupUser

	// Query for group details
	groupQuery := `
		SELECT 
			* 
		FROM 
			public.groups 
		WHERE 
			id = $1
	`
	err := db.db.Get(&group, groupQuery, groupId)
	if err != nil {
		return nil, nil, nil, err
	}

	// Query for group dates
	dateQuery := `
		SELECT 
			* 
		FROM 
			public.group_dates 
		WHERE 
			group_id = $1
	`
	err = db.db.Select(&groupDates, dateQuery, groupId)
	if err != nil {
		return nil, nil, nil, err
	}

	// Query for group users
	userQuery := `
		SELECT 
			* 
		FROM 
			public.group_users 
		WHERE 
			group_id = $1
	`
	err = db.db.Select(&groupUsers, userQuery, groupId)
	if err != nil {
		return nil, nil, nil, err
	}

	return &group, groupDates, groupUsers, nil
}

func (db *Database) GetDates(calendarIDs []uuid.UUID) ([]GroupDate, error) {
	dates := []GroupDate{}
	query := `
		SELECT 
			*
		FROM
			public.group_dates
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
		group_dates (from_date, to_date, group_id)
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

func (db *Database) AddEventToGroup(event GroupEvent) error {
	createdAt := time.Now()
	updatedAt := time.Now()
	startTime := event.StartTime
	endTime := event.EndTime
	eventType := event.Type
	description := event.Description
	title := event.Title
	author := event.Author
	metadata := event.Metadata
	groupID := event.GroupID

	query := `
	INSERT INTO 
		group_events (created_at, updated_at, start_time, end_time, type, description, title, author, metadata, group_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	tx := db.db.MustBegin()
	_, err := tx.Exec(query, createdAt, updatedAt, startTime, endTime, eventType, description, title, author, metadata, groupID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *Database) GetEventFromGroup(groupID uuid.UUID, limit int, offset int) ([]GroupEvent, error) {
	events := []GroupEvent{}
	query := `
        SELECT
            *
        FROM
            public.group_events
        WHERE
            group_id = $1
        LIMIT $2 OFFSET $3
    `

	err := db.db.Select(&events, query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}

	return events, nil
}

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

func (db *Database) AddUsersToGroup(userIDs []uuid.UUID, groupID uuid.UUID) error {
	query := `
		INSERT INTO 
			group_users (user_id, group_id)
		VALUES
			(UNNEST($1::uuid[]), $2);
	`

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, pq.Array(userIDs), groupID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
