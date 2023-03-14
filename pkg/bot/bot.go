package bot

import (
	"context"
	config "subscription-bot/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	bot          *tgbotapi.BotAPI
	client       *mongo.Client
	weatherToken string
	responses    config.Responses
}

func Newbot(bot *tgbotapi.BotAPI, client *mongo.Client, weatherToken string, responses config.Responses) *Bot {
	return &Bot{
		bot:          bot,
		client:       client,
		weatherToken: weatherToken,
		responses:    responses,
	}
}

func (b *Bot) Start() error {
	log.Infof("Authorized on account %s", b.bot.Self.UserName)

	go b.initScheduler()
	log.Info("Scheduler is started")

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
		}

		if update.Message.Location != nil {
			ctx := context.Background()

			if err := b.handleLocationMessage(ctx, update.Message); err != nil {
				return err
			}
		}

		//Handle reply to message with subscriber's time
		if update.Message.ReplyToMessage != nil {
			ctx := context.Background()

			if err := b.handleReplyMessage(ctx, update.Message); err != nil {
				return err
			}
			//Handle other non-valid messages
		} else if !update.Message.IsCommand() && update.Message.Location == nil {
			if err := b.handleOtherMessage(update.Message); err != nil {
				return err
			}
		}
	}
	return nil
}
