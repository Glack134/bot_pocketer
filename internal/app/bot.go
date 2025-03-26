package app

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/polyk005/tg_bot/internal/config"
	"github.com/polyk005/tg_bot/internal/delivery/telegram"
	"github.com/polyk005/tg_bot/internal/repository/inmemory"
	"github.com/polyk005/tg_bot/internal/service"
	"github.com/polyk005/tg_bot/pkg/logger"
)

type Bot struct {
	tgBot   *bot.Bot
	service *service.Service
	logger  logger.Logger
}

func NewBot(cfg *config.Config, log logger.Logger) (*Bot, error) {
	repo := inmemory.New() // вместо postgres.New //postgres.New(context.Background(), cfg.Database.DSN)
	log.Info("Using in-memory repository for testing")
	// if err != nil {
	// 	log.Error("Failed to initialize repository", "error", err)
	// 	return nil, err
	// }

	svc := service.New(repo, log)

	tgBot, err := bot.New(cfg.Telegram.Token, bot.WithDefaultHandler(defaultHandler(log)))
	if err != nil {
		log.Error("Failed to create bot", "error", err)
		return nil, err
	}

	telegram.RegisterHandlers(tgBot, svc, log)

	log.Info("Bot initialized successfully")
	return &Bot{
		tgBot:   tgBot,
		service: svc,
		logger:  log,
	}, nil
}

func defaultHandler(log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			log.Debugw("Received message",
				"chatID", update.Message.Chat.ID,
				"text", update.Message.Text)
		} else if update.CallbackQuery != nil {
			log.Debugw("Received callback",
				"data", update.CallbackQuery.Data)
		}
	}
}

func (b *Bot) Run() error {
	b.logger.Info("Starting bot...")
	b.tgBot.Start(context.Background())
	return nil
}
