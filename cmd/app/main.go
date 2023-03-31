package main

import (
	"subscription-bot/pkg/api"
	b "subscription-bot/pkg/bot"
	"subscription-bot/pkg/config"
	"subscription-bot/pkg/container"
	"subscription-bot/pkg/database"
	lg "subscription-bot/pkg/logger"
	_ "time/tzdata"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	log, err := lg.NewLogger()
	if err != nil {
		return
	}

	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Errorf("failed to load config: %s", err)
		return
	}

	logger, err := lg.LoadLogger(cfg)
	if err != nil {
		log.Errorf("faild to load logger: %s", err)
		return
	}

	container := container.NewBotContainer(cfg, logger)

	weatherStorage := api.NewWeatherStorageService(container)

	requestService := api.NewRequestWeatherService(weatherStorage)

	mongoStorage := database.NewMongoStorageService(container, requestService)

	userStorage := database.NewUserStorageService(mongoStorage)

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Errorf("error creating new BotAPI instance: %s", err)
		return
	}
	botApi.Debug = true

	client, err := userStorage.InitClient(cfg)
	if err != nil {
		logger.Errorf("error creating mongoDB client: %s", err)
		return
	}

	bot := b.Newbot(botApi, userStorage, client, cfg.Response, container, requestService)

	if err := bot.Start(); err != nil {
		logger.Errorf("unable to start bot: %s", err)
		return
	}
}
