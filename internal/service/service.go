package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/polyk005/tg_bot/internal/domain"
	"github.com/polyk005/tg_bot/pkg/logger"
)

type Service struct {
	repo            domain.Repository
	ai              *AIService
	logger          logger.Logger
	schedules       map[int64]map[time.Weekday][]domain.Lesson
	weeklySchedules map[int64]string // Хранит URL фотографий недельных расписаний
	mu              sync.RWMutex
}

func New(repo domain.Repository, ai *AIService, log logger.Logger) *Service {
	return &Service{
		repo:            repo,
		ai:              ai,
		logger:          log,
		schedules:       make(map[int64]map[time.Weekday][]domain.Lesson),
		weeklySchedules: make(map[int64]string),
	}
}

func (s *Service) SaveWeeklyScheduleImage(ctx context.Context, userID int64, imageURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weeklySchedules[userID] = imageURL
	s.logger.Infow("Weekly schedule saved", "userID", userID)
	return nil
}

func (s *Service) GetWeeklyScheduleImage(ctx context.Context, userID int64) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.weeklySchedules[userID]
	return url, exists, nil
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

func (s *Service) SaveSchedule(ctx context.Context, userID int64, day time.Weekday, lessons []domain.Lesson) {
	if _, ok := s.schedules[userID]; !ok {
		s.schedules[userID] = make(map[time.Weekday][]domain.Lesson)
	}
	s.schedules[userID][day] = lessons
}

func (s *Service) GetSchedule(ctx context.Context, userID int64, day time.Weekday) ([]domain.Lesson, error) {
	userSchedule, ok := s.schedules[userID]
	if !ok {
		return nil, fmt.Errorf("schedule not found for user %d", userID)
	}

	lessons, ok := userSchedule[day]
	if !ok {
		return nil, fmt.Errorf("schedule not found for day %v", day)
	}

	return lessons, nil
}

func (s *Service) ProcessAIQuestion(ctx context.Context, question string) (string, error) {
	return s.ai.HandleQuestion(ctx, question)
}

func (s *Service) ProcessScheduleImage(ctx context.Context, userID int64, imageURL string) (string, error) {
	return s.ai.ProcessScheduleImage(ctx, imageURL)
}
