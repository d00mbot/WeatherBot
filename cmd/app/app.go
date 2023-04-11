package app

import (
	"subscription-bot/config"
	"subscription-bot/container"
	"subscription-bot/internal/api"
	"subscription-bot/internal/bot"
	"subscription-bot/internal/database"
	"subscription-bot/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	logger, err := logger.NewBotLogger(cfg)
	if err != nil {
		return err
	}

	responses, err := config.LoadResponses(cfg, logger)
	if err != nil {
		return err
	}

	container := container.NewContainer(cfg, logger)

	forecastService := api.NewForecastService(container)

	weatherService := api.NewWeatherService(forecastService)

	mongoStorage := database.NewMongoStorageService(container, weatherService)

	userStorage := database.NewUserStorageService(mongoStorage)

	client, err := database.NewMongoClient(container)
	if err != nil {
		return err
	}

	botApi, err := tgbotapi.NewBotAPI(container.GetConfig().TelegramToken)
	if err != nil {
		logger.Errorf("error creating new BotAPI:\n'%v'", err)
		return err
	}
	botApi.Debug = true

	bot := bot.Newbot(botApi, client, userStorage, weatherService, container, responses)

	if err := bot.Start(); err != nil {
		logger.Errorf("unable to start bot:\n'%v'", err)
		return err
	}

	return nil
}
