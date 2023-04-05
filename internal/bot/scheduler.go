package bot

import (
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const timeLayout = "15:04"

func (b *Bot) startScheduler() {
	ns := gocron.NewScheduler(time.UTC)

	ns.Every(1).Hour().Do(func() {
		b.sendScheduledMessage()
	})
	ns.StartAsync()
}

func (b *Bot) sendScheduledMessage() error {
	logger := b.container.GetLogger()

	users, err := b.storage.FindAll(b.client)
	if err != nil {
		logger.Errorf("faild to get all users: %s", err)
		return err
	}

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

		ns.Every(1).Day().At(userTime).Do(func() {
			if err := b.sendForecast(u.Latitude, u.Longitude, u.ChatID); err != nil {
				return
			}
		})
		ns.StartAsync()
	}

	return nil
}

func (b *Bot) sendForecast(lat float64, lon float64, chatID int64) error {
	logger := b.container.GetLogger()

	forecast, _, err := b.weatherService.GetForecast(lat, lon)
	if err != nil {
		logger.Errorf("faild to create weather forecast message: %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, forecast)

	_, err = b.bot.Send(msg)
	if err != nil {
		logger.Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}
