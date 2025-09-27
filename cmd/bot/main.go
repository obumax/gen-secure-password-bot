package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"gen-secure-password-bot/internal/i18n"
	"gen-secure-password-bot/internal/session"
	"gen-secure-password-bot/internal/telegram"
)

func main() {

	// Загрузка окружения из .env / Load .env file

	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found: %v", err)
	}

	// Инициализация i18n / Initialize i18n

	if err := i18n.InitBundle(); err != nil {
		log.Fatalf("i18n init failed: %v", err)
	}

	// Читается REDIS_URL и создается Store / Read REDIS URL and create Store

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not set")
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}

	// Создается и инициализируется Store для сессий / Create and initialize Store for sessions

	store := session.NewRedisStore(opts)
	session.InitStore(store)

	// Читается токен бота / Read bot token

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	// Бот создается и запускается / Bot is created and strarted

	bot, err := telegram.NewBot(token, store)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}
	if err := bot.Run(context.Background()); err != nil {
		log.Fatalf("bot stopped with error: %v", err)
	}
}
