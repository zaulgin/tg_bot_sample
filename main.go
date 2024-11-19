package main

import (
	"flag"
	"log"

	tgClient "github.com/zaulgin/tg_bot_sample/clients/telegram"
	event_consumer "github.com/zaulgin/tg_bot_sample/consumer/event-consumer"
	"github.com/zaulgin/tg_bot_sample/events/telegram"
	"github.com/zaulgin/tg_bot_sample/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProccesor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProccesor, eventsProccesor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
