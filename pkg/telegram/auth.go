package telegram

import (
	"context"
	"fmt"
	"tg_bot/pkg/storage"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.createAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf(b.messages.Responses.Start, authLink)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err = b.Bot.Send(msg)

	return err
}

func (b *Bot) createAuthorizationLink(chatID int64) (string, error) {
	redirectUrl := b.generateRedirectURL(chatID)
	token, err := b.client.GetRequestToken(context.Background(), b.redirectUrl)
	if err != nil {
		return "", err
	}

	if err := b.storage.Save(chatID, token, storage.RequestTokens); err != nil {
		return "", err
	}

	return b.client.GetAuthorezationURL(token, redirectUrl)
}

func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatID)
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.storage.Get(chatID, storage.AccessTokens)
}
