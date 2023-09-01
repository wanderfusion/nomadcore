package calendar

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/jwt"
	"github.com/akxcix/nomadcore/pkg/repositories/calendar"
)

type Service struct {
	JwtManager *jwt.JwtManager
	CalRepo    *calendar.Database
}

func New(dbConf *config.DatabaseConfig, jwtConf *config.Jwt) *Service {
	if dbConf == nil {
		log.Fatal().Msg("dbConf is nil")
	}

	if jwtConf == nil {
		log.Fatal().Msg("jwtConf is nil")
	}

	calRepo := calendar.New(dbConf)
	jwtManager := jwt.New(jwtConf.Secret, jwtConf.ValidMins)

	svc := &Service{
		JwtManager: jwtManager,
		CalRepo:    calRepo,
	}

	return svc
}

func (s *Service) CreateCalendar(userId uuid.UUID, name, visibility string) (string, error) {
	err := s.CalRepo.CreateCalendar(userId, name, visibility)
	if err != nil {
		return "", err
	}

	msg := "Successfully added calendar"
	return msg, nil
}

func (s *Service) ValidateJwt(token string) (*jwt.Claims, bool) {
	claims, err := s.JwtManager.Verify(token)
	if err != nil {
		return nil, false
	}

	return claims, true
}
