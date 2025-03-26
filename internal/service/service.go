package service

import (
	"context"

	"github.com/polyk005/tg_bot/internal/domain"
	"github.com/polyk005/tg_bot/pkg/logger"
)

type Service struct {
	repo   domain.Repository
	logger logger.Logger
}

func New(repo domain.Repository, log logger.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: log,
	}
}

func (s *Service) ProcessStartCommand(ctx context.Context, userID int64) error {
	s.logger.Infow("Processing start command", "userID", userID)

	exists, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}

	if !exists {
		return s.repo.CreateUser(ctx, userID)
	}

	return nil
}
