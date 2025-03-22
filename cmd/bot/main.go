package main

import (
	"log"
	"tg_bot/pkg/config"
	"tg_bot/pkg/telegram"

	"github.com/boltdb/bolt"
	"github.com/containerd/containerd/snapshots/storage"
	"github.com/docker/docker/api/server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initBolt()
	if err != nil {
		log.Fatal(err)
	}
	storage := boltdbNewTokenStorage(db)

	bot := telegram.NewBot(botApi, pocketClient, cfg.AuthServerURL, storage, cfg.Messages)

	redirectServer := server.NewAuthServer(cfg.BotURL, storage, pocketClient)

	go func() {
		if err := redirectServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initBolt() (*bolt.DB, error) {
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.RequestToken))
		return err
	}); err != nil {
		return nil, err
	}
	return db, err
}
