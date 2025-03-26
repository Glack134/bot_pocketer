package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/polyk005/tg_bot/internal/service"
	"github.com/polyk005/tg_bot/pkg/logger"
)

func RegisterHandlers(b *bot.Bot, svc *service.Service, log logger.Logger) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler(svc, log))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler(log))
}

func startHandler(svc *service.Service, log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userID := update.Message.From.ID

		if err := svc.ProcessStartCommand(ctx, userID); err != nil {
			log.Errorw("Failed to process start command", "error", err)
			return
		}

		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Welcome to the bot!",
		})

		if err != nil {
			log.Errorw("Failed to send message", "error", err)
		}
	}
}

func helpHandler(log logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Available commands:\n/start - Start bot\n/help - Show help",
		})

		if err != nil {
			log.Errorw("Failed to send help message", "error", err)
		}
	}
}
