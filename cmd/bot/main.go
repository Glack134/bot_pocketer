package main

import (
	"fmt"
	"log"

	"github.com/polyk005/tg_bot/internal/app"
	"github.com/polyk005/tg_bot/internal/config"
	"github.com/polyk005/tg_bot/pkg/logger"
)

func main() {
	fmt.Println("Starting Telegram Bot...")

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.New(cfg.LogLevel)
	defer log.Sync()

	bot, err := app.NewBot(cfg, log)
	if err != nil {
		log.Fatal("Failed to create bot", "error", err)
	}

	if err := bot.Run(); err != nil {
		log.Fatal("Bot stopped with error", "error", err)
	}
}
