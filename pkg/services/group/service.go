package group

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/akxcix/nomadcore/pkg/clients/passport"
	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/repositories/group"
	"github.com/akxcix/nomadcore/pkg/services"
)

type Service struct {
	groupRepo      *group.Database
	passportClient *passport.Client
}

func New(dbConf *config.DatabaseConfig, passportClientConf *config.PassportClient) *Service {
	if dbConf == nil {
		log.Fatal().Msg("dbConf is nil")
	}

	groupRepo := group.New(dbConf)
	passportClient := passport.NewClient(passportClientConf.Host)

	svc := &Service{
		groupRepo:      groupRepo,
		passportClient: passportClient,
	}

	return svc
}

func (s *Service) CreateGroup(userId uuid.UUID, name, description string) (string, services.ServiceError) {
	groupID, err := s.groupRepo.CreateGroup(userId, name, description)
	if err != nil {
		log.Error().Err(err).Msg("unable to add calendar to DB")
		return "", ErrFailedDBWrite
	}

	err = s.groupRepo.AddUsersToGroup([]uuid.UUID{userId}, groupID)
	if err != nil {
		log.Error().Err(err).Msg("unable to add user to DB")
		return "", ErrFailedDBWrite
	}

	msg := "Successfully added calendar"
	return msg, nil
}

func (s *Service) GetGroups(userID uuid.UUID) ([]group.Group, services.ServiceError) {
	groups, err := s.groupRepo.GetGroups(userID)
	if err != nil {
		log.Error().Err(err).Msg("unable to get groups from DB")
		return nil, ErrFailedDBRead
	}

	return groups, nil
}

func (s *Service) GetGroupDetails(userID, groupId uuid.UUID) (*group.Group, []group.GroupDate, []group.GroupUser, services.ServiceError) {
	groupDetails, groupDates, groupUsers, err := s.groupRepo.GetGroupWithDetails(groupId)
	groupUsers = append(groupUsers, group.GroupUser{UserID: userID})
	if err != nil {
		log.Error().Err(err).Msg("Unable to get group details from DB")
		return nil, nil, nil, ErrFailedDBRead
	}

	userAuthorised := false
	for _, user := range groupUsers {
		if user.UserID == userID {
			userAuthorised = true
			break
		}
	}

	if !userAuthorised {
		log.Error().Msg("User not authorised to view group")
		return nil, nil, nil, ErrUserForbidden
	}

	return groupDetails, groupDates, groupUsers, nil
}

func (s *Service) GetDates(calendarIDs []uuid.UUID) ([]group.GroupDate, services.ServiceError) {
	dates, err := s.groupRepo.GetDates(calendarIDs)
	if err != nil {
		log.Error().Err(err).Msg("unable to get dates from DB")
		return nil, ErrFailedDBRead
	}

	return dates, nil
}

func (s *Service) AddDatesToGroup(userID, groupID uuid.UUID, dates Dates) (string, services.ServiceError) {
	if dates.To.Before(dates.From) {
		log.Info().Msg("invalid dates")
		return "", ErrInvalidRequest
	}
	err := s.groupRepo.AddDatesToGroup(userID, groupID, dates.From, dates.To)
	if err != nil {
		log.Error().Err(err).Msg("unable to add dates to DB")
		return "", ErrFailedDBWrite
	}
	return "successfully added", nil
}

func (s *Service) AddUsersToGroup(usernames []string, groupID uuid.UUID) (string, services.ServiceError) {
	if usernames == nil {
		log.Info().Msg("invalid usernames")
		return "", ErrInvalidRequest
	}
	usernameMap, err := s.getUserIDsFromUsernames(usernames)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user IDs from passport")
		return "", ErrFailedClientCall
	}
	userIDs := make([]uuid.UUID, 0)
	for _, userID := range usernameMap {
		userIDs = append(userIDs, userID)
	}
	err = s.groupRepo.AddUsersToGroup(userIDs, groupID)
	if err != nil {
		log.Error().Err(err).Msg("unable to add dates to DB")
		return "", ErrFailedDBWrite
	}
	return "successfully added", nil
}

func (s *Service) getUserIDsFromUsernames(usernames []string) (map[string]uuid.UUID, error) {
	userIDs := make(map[string]uuid.UUID)
	res, err := s.passportClient.CachedGetUsersFromUsernames(usernames, 5*time.Hour)
	if err != nil {
		log.Error().Err(err).Msg("unable to get users from passport")
		return nil, errors.New("unable to get users from passport")
	}
	if len(res.Data) == 0 {
		log.Error().Msg("no users found")
		return nil, errors.New("no users found")
	}
	for _, user := range res.Data {
		userIDs[user.Username] = user.ID
	}

	return userIDs, nil
}

func (s *Service) GetUserProfileByID(userID uuid.UUID) (group.UserProfile, error) {
	userProfile, err := s.groupRepo.GetUserProfileByID(userID)
	if err != nil {
		return userProfile, err
	}

	return userProfile, nil
}
