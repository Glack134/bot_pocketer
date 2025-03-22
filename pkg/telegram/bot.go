package telegram

import (
	"tg_bot/pkg/config"
	"tg_bot/pkg/storage"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	client      *pocket.Client
	redirectURL string

	storage storage.TokenStorage

	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, redirectURL string, storage storage.TokenStorage, messages config.Messages) *Bot {
	return &Bot{
		bot:         bot,
		client:      client,
		redirectURL: redirectURL,
		storage:     storage,
		messages:    messages,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate()
	u.Timeout = 60

	updates, err := b.bot.GetUpdateChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
	return nil
}
