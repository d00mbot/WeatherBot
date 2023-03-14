package bot

import (
	api "subscription-bot/pkg/api"
	database "subscription-bot/pkg/database"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) initScheduler() {
	ns := gocron.NewScheduler(time.UTC)

	ns.Every(1).Hour().Do(func() {
		b.sendScheduledMessage()
	})
	ns.StartAsync()
}

func (b *Bot) sendScheduledMessage() error {
	subs, err := database.GetAllSubscribers(b.client)
	if err != nil {
		log.Errorf("error getting all subscribres: %s", err)
		return err
	}

	for _, sub := range subs {
		s := sub

		userLocation, err := time.LoadLocation(s.Timezone)
		if err != nil {
			log.Errorf("error loading location: %s", err)
			return err
		}

		ns := gocron.NewScheduler(userLocation)

		userTime, err := time.Parse("15:04", s.Time)
		if err != nil {
			log.Errorf("error parsing user time from mongoDB: %s", err)
			return err
		}

		ns.Every(1).Day().At(userTime).Do(func() {
			if err := b.sendWeatherMessage(s.Latitude, s.Longitude, s.ChatID); err != nil {
				return
			}
		})
		ns.StartAsync()
	}

	return nil
}

func (b *Bot) sendWeatherMessage(lat float64, long float64, chatID int64) error {
	weatherMsg, _, err := api.RequestWeatherForecast(lat, long, b.weatherToken)
	if err != nil {
		log.Errorf("error sending weather forecast message: %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, weatherMsg)

	if _, err = b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}
