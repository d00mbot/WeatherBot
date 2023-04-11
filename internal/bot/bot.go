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

type bot struct {
	bot            *tgbotapi.BotAPI
	client         *mongo.Client
	storage        database.UserStorageService
	weatherService api.WeatherService
	container      container.Container
	responses      *models.Responses
}

type Bot interface {
	Start() error
}

func Newbot(
	botAPI *tgbotapi.BotAPI,
	client *mongo.Client,
	storage database.UserStorageService,
	weatherService api.WeatherService,
	container container.Container,
	responses *models.Responses,
) Bot {
	return &bot{
		bot:            botAPI,
		client:         client,
		storage:        storage,
		weatherService: weatherService,
		container:      container,
		responses:      responses,
	}
}

func (b *bot) Start() error {
	b.container.GetLogger().Infof("Authorized on account %s", b.bot.Self.UserName)

	go func() {
		if err := b.startScheduler(); err != nil {
			b.container.GetLogger().Errorf("unable to start scheduler:\n%v", err)
		}
	}()
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
