package bot

import (
	"context"
	"subscription-bot/container"
	"subscription-bot/internal/api"
	"subscription-bot/internal/database"
	"subscription-bot/internal/models"

	"go.mongodb.org/mongo-driver/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	storage        database.UserStorageService
	client         *mongo.Client
	responses      *models.Responses
	container      container.BotContainer
	weatherService api.WeatherService
}

func Newbot(
	bot *tgbotapi.BotAPI,
	storage database.UserStorageService,
	client *mongo.Client,
	responses *models.Responses,
	container container.BotContainer,
	weatherService api.WeatherService) *Bot {

	return &Bot{
		bot:            bot,
		storage:        storage,
		client:         client,
		responses:      responses,
		container:      container,
		weatherService: weatherService,
	}
}

func (b *Bot) Start() error {
	b.container.GetLogger().Infof("Authorized on account %s", b.bot.Self.UserName)

	go b.startScheduler()
	b.container.GetLogger().Info("Scheduler is started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				return err
			}
		} else {
			ctx := context.Background()

			if err := b.handleMessage(ctx, update.Message); err != nil {
				return err
			}
		}
	}

	return nil
}
