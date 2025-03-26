package telegram

import (
	"context"

	"github.com/Glack134/websocket/pkg/service"
	"github.com/docker/docker/daemon/logger"
)

func RegisterHandlers(b *bot.Bot, svc *service.Service, logger logger.Logger) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler(svc, logger))
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler(svc, logger))
	// Добавьте другие обработчики здесь
}

func startHandler(svc *service.Service, logger logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Добро пожаловать! Я ваш бот на Go.",
		})
		if err != nil {
			logger.Error("failed to send message", "error", err)
		}
	}
}

func helpHandler(svc *service.Service, logger logger.Logger) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Доступные команды:\n/start - начать работу\n/help - помощь",
		})
		if err != nil {
			logger.Error("failed to send message", "error", err)
		}
	}
}
