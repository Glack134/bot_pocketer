package telegram

import "errors"

var (
	invalidUrlError   = errors.New("url is valid")
	unableToSaveError = errors.New("unable to save link to Pocket")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string

	switch err {
	case invalidUrlError:
		messageText = b.messages.Errors.InvalidURL
	case unableToSaveError:
		messageText = b.messages.Errors.UnableToSave
	default:
		messageText = b.messages.Errors.Default
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	b.bot.Send(msg)
}
