package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// handleUpdate обрабатывает входящие обновления / handleUpdate processes incoming updates

func (b *Bot) handleUpdate(u tgbotapi.Update) {
	if u.CallbackQuery != nil {
		b.handleCallback(u.CallbackQuery)
	} else if u.Message != nil {
		b.handleMessage(u.Message)
	}
}
