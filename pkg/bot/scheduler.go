package bot

import (
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const timeLayout = "15:04"

func (b *Bot) initScheduler() {
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
		logger.Errorf("error getting all users: %s", err)
		return err
	}

	for _, user := range users {
		u := user

		userLocation, err := time.LoadLocation(u.Timezone)
		if err != nil {
			logger.Errorf("error loading location: %s", err)
			return err
		}

		ns := gocron.NewScheduler(userLocation)

		userTime, err := time.Parse(timeLayout, u.Time)
		if err != nil {
			logger.Errorf("error parsing time: %s", err)
			return err
		}

		ns.Every(1).Day().At(userTime).Do(func() {
			if err := b.sendWeatherMessage(u.Latitude, u.Longitude, u.ChatID); err != nil {
				return
			}
		})
		ns.StartAsync()
	}

	return nil
}

func (b *Bot) sendWeatherMessage(lat float64, lon float64, chatID int64) error {
	logger := b.container.GetLogger()

	weatherMsg, _, err := b.request.RequestWeather(lat, lon, b.container.GetConfig())
	if err != nil {
		logger.Errorf("error creating weather forecast message: %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, weatherMsg)

	if _, err = b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}
