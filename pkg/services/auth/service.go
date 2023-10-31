package auth

import (
	"github.com/rs/zerolog/log"

	"github.com/wanderfusion/nomadcore/pkg/config"
	"github.com/wanderfusion/nomadcore/pkg/jwt"
)

type Service struct {
	jwtManager *jwt.JwtManager
}

func New(jwtConf *config.Jwt) *Service {
	if jwtConf == nil {
		log.Fatal().Msg("jwtConf is nil")
	}

	jwtManager := jwt.New(jwtConf.Secret, jwtConf.ValidMins)

	svc := &Service{
		jwtManager: jwtManager,
	}

	return svc
}

func (s *Service) ValidateJwt(token string) (*jwt.Claims, bool) {
	claims, err := s.jwtManager.Verify(token)
	if err != nil {
		return nil, false
	}

	return claims, true
}
