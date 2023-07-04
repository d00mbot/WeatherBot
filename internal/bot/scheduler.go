package bot

import (
	"subscription-bot/internal/models"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const timeLayout = "15:04"

func (b *bot) startScheduler() error {
	ns := gocron.NewScheduler(time.UTC)

	_, err := ns.Every(1).Hour().Do(
		func() {
			users, err := b.storage.FindAll(b.client)
			if err != nil {
				b.container.GetLogger().Errorf("start scheduler:\n'%v'", err)
			}

			if err := b.sendScheduledMessage(users); err != nil {
				b.container.GetLogger().Errorf("start scheduler:\n'%v'", err)
			}
		},
	)
	if err != nil {
		b.container.GetLogger().Errorf("start scheduler:\n'%v'", err)
		return err
	}

	ns.StartAsync()

	return nil
}

func (b *bot) sendScheduledMessage(users []models.User) error {
	logger := b.container.GetLogger()

	for _, user := range users {
		u := user

		userLocation, err := time.LoadLocation(u.Timezone)
		if err != nil {
			logger.Errorf("faild to load location: %s", err)
			return err
		}

		ns := gocron.NewScheduler(userLocation)

		userTime, err := time.Parse(timeLayout, u.Time)
		if err != nil {
			logger.Errorf("faild to parse time: %s", err)
			return err
		}

		_, err = ns.Every(1).Day().At(userTime).Do(
			func() {
				if err := b.sendForecast(u.Latitude, u.Longitude, u.ChatID); err != nil {
					logger.Errorf("unable to send forecast:\n'%v'", err)
				}
			},
		)
		if err != nil {
			logger.Errorf("unable to send scheduled message:\n'%v'", err)
			return err
		}

		ns.StartAsync()
	}

	return nil
}

func (b *bot) sendForecast(lat float64, lon float64, chatID int64) error {
	forecast, _, err := b.weatherService.GetForecast(lat, lon)
	if err != nil {
		b.container.GetLogger().Errorf("error creating forecast message: %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, forecast)

	_, err = b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}
