package cmd

import (
	"fmt"
	"log"

	"github.com/containerd/containerd/cmd/ctr/app"
)

func main() {
	fmt.Println("Start Bot to @tg_bot_chirik_10")

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.LogLevel)

	bot, err := app.NewBot(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create bot", "error", err)
	}

	if err := bot.Run(); err != nil {
		logger.Fatal("Bot stopped with error", "error", err)
	}
}
