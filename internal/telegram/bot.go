package telegram

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gen-secure-password-bot/internal/session"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	store session.Store
}

// NewBot создает нового бота с заданным токеном и хранилищем сессий / NewBot creates a new bot with the given token and session store

func NewBot(token string, store session.Store) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", api.Self.UserName)
	return &Bot{api: api, store: store}, nil
}

// Run запускает бота и обрабатывает обновления / Run starts the bot and processes updates

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case upd := <-updates:
			b.handleUpdate(upd)
		}
	}
}
