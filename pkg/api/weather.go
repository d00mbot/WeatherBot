package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"subscription-bot/pkg/config"
	"subscription-bot/pkg/container"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WeatherStorageService struct {
	container container.BotContainer
}

func NewWeatherStorageService(c container.BotContainer) WeatherStorageService {
	return WeatherStorageService{container: c}
}

type Weather struct {
	TimeZone string `json:"timezone"`
	Daily    []struct {
		Temperature struct {
			Day   float32 `json:"day"`
			Min   float32 `json:"min"`
			Max   float32 `json:"max"`
			Night float32 `json:"night"`
		} `json:"temp"`
		FeelsLike struct {
			Day   float32 `json:"day"`
			Night float32 `json:"night"`
		} `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Clouds    int     `json:"clouds"`
		WindSpeed float32 `json:"wind_speed"`
		Weather   []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"daily"`
}

func (w *WeatherStorageService) requestWeatherForecast(lat float64, lon float64, cfg *config.Config) (string, string, error) {
	logger := w.container.GetLogger()

	if lat == 0.0 {
		err := errors.New("latitude is empty")
		logger.Errorf("faild to request weather forecast: %s", err)
		return "", "", err
	}

	if lon == 0.0 {
		err := errors.New("longitude is empty")
		logger.Errorf("faild to request weather forecast: %s", err)
		return "", "", err
	}

	response, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=current,minutely,hourly,alerts&units=metric&appid=%s",
		lat,
		lon,
		cfg.WeatherToken))
	if err != nil {
		logger.Errorf("error getting weather api response: %s", err)
		return "", "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Errorf("error reading weather api response: %s", err)
		return "", "", err
	}

	var resp Weather

	err = json.Unmarshal(body, &resp)
	if err != nil {
		logger.Errorf("error unmarshaling weather api response: %s", err)
		return "", "", err
	}

	return string(fmt.Sprintf("Weather forecast for the day:\nTemperature: min %dC, max %dC, Day %dC, Night %dC\nFeels like: Day %dC, Night %dC\nHumidity: %d%%\nCloudiness: %d%%\nWind speed: %.2fm/s\nWeather: %s(%s)",
		int(resp.Daily[0].Temperature.Min), int(resp.Daily[0].Temperature.Max), int(resp.Daily[0].Temperature.Day), int(resp.Daily[0].Temperature.Night),
		int(resp.Daily[0].FeelsLike.Day), int(resp.Daily[0].FeelsLike.Night), resp.Daily[0].Humidity, resp.Daily[0].Clouds, resp.Daily[0].WindSpeed,
		resp.Daily[0].Weather[0].Main, resp.Daily[0].Weather[0].Description)), resp.TimeZone, nil
}

func (w *WeatherStorageService) requestUserTimezone(message *tgbotapi.Message, cfg *config.Config) (string, error) {
	logger := w.container.GetLogger()

	_, timeZone, err := w.requestWeatherForecast(message.Location.Latitude, message.Location.Longitude, cfg)
	if err != nil {
		logger.Errorf("error getting user timezone: %s", err)
		return "", err
	}
	return timeZone, nil
}
