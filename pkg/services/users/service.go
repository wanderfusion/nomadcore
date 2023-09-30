package users

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/akxcix/nomadcore/pkg/clients/passport"
	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/repositories/users"
	"github.com/akxcix/nomadcore/pkg/services"
)

type Service struct {
	userRepo       *users.Database
	passportClient *passport.Client
}

func New(dbConf *config.DatabaseConfig, passportClientConf *config.PassportClient) *Service {
	if dbConf == nil {
		log.Fatal().Msg("dbConf is nil")
	}

	userRepo := users.New(dbConf)
	passportClient := passport.NewClient(passportClientConf.Host)

	svc := &Service{
		userRepo:       userRepo,
		passportClient: passportClient,
	}

	return svc
}

func (s *Service) GetUserProfile(username string) (users.UserProfile, services.ServiceError) {
	usernameMap, err := s.getUserIDsFromUsernames([]string{username})
	if err != nil {
		log.Error().Err(err).Msg("unable to get user ID from passport")
		return users.UserProfile{}, ErrFailedClientCall
	}
	userIDs := make([]uuid.UUID, 0)
	for _, userID := range usernameMap {
		userIDs = append(userIDs, userID)
	}
	if len(userIDs) != 1 {
		log.Error().Err(err).Msg("expected the length to be 1")
		return users.UserProfile{}, ErrInvalidRequest
	}

	userID := userIDs[0]

	userProfile, err := s.userRepo.GetUserProfileByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msg("no row found")
			return users.UserProfile{}, ErrNoRowFound // Assuming ErrNoRowFound is a predefined error constant
		}
		log.Error().Err(err).Msg("something went wrong")
		return users.UserProfile{}, ErrFailedDBRead
	}

	return userProfile, nil
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
