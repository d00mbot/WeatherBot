package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"subscription-bot/container"
	"subscription-bot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeatherServiceStorage struct {
	container container.BotContainer
}

func NewWeatherServiceStorage(c container.BotContainer) WeatherServiceStorage {
	return WeatherServiceStorage{container: c}
}

func (w *WeatherServiceStorage) getWeatherForecast(lat float64, lon float64) (string, string, error) {
	logger := w.container.GetLogger()

	if err := validateGeopoints(lat, lon); err != nil {
		logger.Errorf("faild to get weather forecast: %s", err)
		return "", "", err
	}

	response, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=current,minutely,hourly,alerts&units=metric&appid=%s",
		lat,
		lon,
		w.container.GetConfig().WeatherToken))
	if err != nil {
		logger.Errorf("faild to get weather api response: %s", err)
		return "", "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("faild to read weather api response: %s", err)
		return "", "", err
	}

	var resp models.Weather

	if err := json.Unmarshal(body, &resp); err != nil {
		logger.Errorf("faild to unmarshal weather api response: %s", err)
		return "", "", err
	}

	return generateForecast(resp), resp.TimeZone, nil
}

func validateGeopoints(lat float64, lon float64) error {
	if lat == 0.0 {
		err := errors.New("latitude is empty")
		return err
	} else if lon == 0.0 {
		err := errors.New("longitude is empty")
		return err
	}

	return nil
}

func generateForecast(w models.Weather) string {
	forecast := fmt.Sprintf("Forecast for the day:\n\nWeather: %s(%s)\n\nTemperature:\nDay: %dC\nNight: %dC\nMin: %dC\nMax: %dC\n\nFeels like:\nDay: %dC\nNight: %dC\n\nHumidity: %d%%\nCloudiness: %d%%\nWind speed: %.2fm/s.",
		w.Daily[0].Weather[0].Main,
		w.Daily[0].Weather[0].Description,
		int(w.Daily[0].Temperature.Day),
		int(w.Daily[0].Temperature.Night),
		int(w.Daily[0].Temperature.Min),
		int(w.Daily[0].Temperature.Max),
		int(w.Daily[0].FeelsLike.Day),
		int(w.Daily[0].FeelsLike.Night),
		w.Daily[0].Humidity,
		w.Daily[0].Clouds,
		w.Daily[0].WindSpeed,
	)

	return forecast
}

func (w *WeatherServiceStorage) getUserTimezone(message *tgbotapi.Message) (string, error) {
	_, timezone, err := w.getWeatherForecast(message.Location.Latitude, message.Location.Longitude)
	if err != nil {
		w.container.GetLogger().Errorf("faild to get user's timezone: %s", err)
		return "", err
	}

	return timezone, nil
}
