package api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type weatherService struct {
	service WeatherServiceStorage
}

type WeatherService interface {
	GetForecast(lat float64, lon float64) (string, string, error)
	GetTimezone(message *tgbotapi.Message) (string, error)
}

func NewWeatherService(ws WeatherServiceStorage) WeatherService {
	return &weatherService{service: ws}
}

func (w *weatherService) GetForecast(lat float64, lon float64) (string, string, error) {
	forecast, timezone, err := w.service.getWeatherForecast(lat, lon)
	if err != nil {
		return "", "", err
	}

	return forecast, timezone, nil
}

func (w *weatherService) GetTimezone(message *tgbotapi.Message) (string, error) {
	timezone, err := w.service.getUserTimezone(message)
	if err != nil {
		return "", err
	}

	return timezone, nil
}