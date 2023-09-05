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

func (s *Service) CreateCalendar(userId uuid.UUID, name, visibility string) (string, services.ServiceError) {
	err := s.calRepo.CreateCalendar(userId, name, visibility)
	if err != nil {
		log.Error().Err(err).Msg("unable to add calendar to DB")
		return "", ErrFailedDBWrite
	}

	msg := "Successfully added calendar"
	return msg, nil
}

func (s *Service) GetCalendars(userID uuid.UUID, visibility string) ([]calendar.Calendar, services.ServiceError) {
	calendars, err := s.calRepo.GetCalendars(userID, calendar.Visibility(visibility))
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendars from DB")
		return nil, ErrFailedDBRead
	}

	return calendars, nil
}

func (s *Service) GetDates(calendarIDs []uuid.UUID) ([]calendar.Date, services.ServiceError) {
	dates, err := s.calRepo.GetDates(calendarIDs)
	if err != nil {
		log.Error().Err(err).Msg("unable to get dates from DB")
		return nil, ErrFailedDBRead
	}

	return dates, nil
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
