package api

import (
	"subscription-bot/pkg/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type requestWeatherService struct {
	request WeatherStorageService
}

type RequestWeatherService interface {
	RequestWeather(lat float64, lon float64, cfg *config.Config) (string, string, error)
	RequestTimezone(message *tgbotapi.Message, cfg *config.Config) (string, error)
}

func NewRequestWeatherService(w WeatherStorageService) RequestWeatherService {
	return &requestWeatherService{request: w}
}

func (r *requestWeatherService) RequestWeather(lat float64, lon float64, cfg *config.Config) (string, string, error) {
	weather, timezone, err := r.request.requestWeatherForecast(lat, lon, cfg)
	if err != nil {
		return "", "", err
	}

	return weather, timezone, nil
}

func (r *requestWeatherService) RequestTimezone(message *tgbotapi.Message, cfg *config.Config) (string, error) {
	timezone, err := r.request.requestUserTimezone(message, cfg)
	if err != nil {
		return "", err
	}

	return timezone, nil
}
