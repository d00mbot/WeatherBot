package main

import (
	"subscription-bot/config"
	"subscription-bot/container"
	"subscription-bot/internal/api"
	"subscription-bot/internal/bot"
	"subscription-bot/internal/database"
	"subscription-bot/logger"
	_ "time/tzdata"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	// log, err := logger.NewLogger()
	// if err != nil {
	// 	return
	// }

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("failed to load config: %s", err)
		return
	}

	responses, err := config.LoadResponses(cfg)
	if err != nil {
		return
	}

	logger, err := logger.LoadLogger(cfg)
	if err != nil {
		log.Errorf("faild to load logger: %s", err)
		return
	}

	container := container.NewBotContainer(cfg, logger)

	weatherStorage := api.NewWeatherStorageService(container)

	weatherService := api.NewWeatherService(weatherStorage)

	mongoStorage := database.NewMongoStorageService(container, weatherService)

	userStorage := database.NewUserStorageService(mongoStorage)

	botApi, err := tgbotapi.NewBotAPI(container.GetConfig().TelegramToken)
	if err != nil {
		logger.Errorf("faild to create new BotAPI instance: %s", err)
		return
	}
	botApi.Debug = true

	client, err := database.NewMongoClient(container)
	if err != nil {
		return
	}

	bot := bot.Newbot(botApi, userStorage, client, responses, container, weatherService)

	if err := bot.Start(); err != nil {
		logger.Errorf("unable to start bot: %s", err)
		return
	}
}
