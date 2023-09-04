package calendar

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/akxcix/nomadcore/pkg/config"
	"github.com/akxcix/nomadcore/pkg/repositories/calendar"
	"github.com/akxcix/nomadcore/pkg/services"
)

type Service struct {
	calRepo *calendar.Database
}

func New(dbConf *config.DatabaseConfig) *Service {
	if dbConf == nil {
		log.Fatal().Msg("dbConf is nil")
	}

	calRepo := calendar.New(dbConf)

	svc := &Service{
		calRepo: calRepo,
	}

	return svc
}

func (s *Service) CreateCalendar(userId uuid.UUID, name, visibility string) (string, error) {
	err := s.calRepo.CreateCalendar(userId, name, visibility)
	if err != nil {
		return "", err
	}

	msg := "Successfully added calendar"
	return msg, nil
}

func (s *Service) GetCalendars(userID uuid.UUID, visibility string) ([]calendar.Calendar, error) {
	calendars, err := s.calRepo.GetCalendars(userID, calendar.Visibility(visibility))
	if err != nil {
		return nil, err
	}

	return calendars, nil
}

func (s *Service) AddDatesToCalendar(userID, calendarID uuid.UUID, dates Dates) (string, services.ServiceError) {
	if dates.To.Before(dates.From) {
		log.Info().Msg("invalid dates")
		return "", ErrInvalidRequest
	}
	err := s.calRepo.AddDatesToCalendar(userID, calendarID, dates.From, dates.To)
	if err != nil {
		log.Error().Err(err).Msg("unable to add dates to DB")
		return "", ErrFailedDBWrite
	}
	return "successfully added", nil
}
