package app

import (
	"context"

	"github.com/docker/docker/daemon/logger"
	"github.com/polyk005/tg_bot/internal/config"
)

type Bot struct {
	tgBot   *bot.Bot
	service *service.Service
	logger  *logger.Logger
}

func NewBot(cfg *config.Config, logger logger.Logger) (*Bot, error) {
	repo, err := postgres.New(context.Background(), cfg.Database.DSN)
	if err != nil {
		return nil, err
	}

	svc := service.New(repo, logger)

	tgBot, err := bot.New(cfg.Telegram.Token)
	if err != nil {
		return nil, err
	}

	telegram.RegisterHandler(tgBot, svc, logger)

	return &Bot{
		tgBot:   tgBot,
		service: svc,
		logger:  logger,
	}, nil
}

func (b *Bot) Run() error {
	b.logger.Info("Starting Bot")
	b.tgBot.Start(context.Background())
	return nil
}
