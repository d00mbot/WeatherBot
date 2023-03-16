package main

import (
	b "subscription-bot/pkg/bot"
	config "subscription-bot/pkg/config"
	database "subscription-bot/pkg/database"
	_ "time/tzdata"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		return
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Errorf("error creating new BotAPI instance: %s", err)
		return
	}
	botApi.Debug = true

	client, err := database.InitMongoClient(cfg.MongoURI)
	if err != nil {
		log.Errorf("error creating mongoDB client: %s", err)
		return
	}

	bot := b.Newbot(botApi, client, cfg.WeatherToken, cfg.Responses)

	if err := bot.Start(); err != nil {
		log.Errorf("unable to start bot: %s", err)
		return
	}
}
