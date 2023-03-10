package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

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

func RequestWeatherForecast(lat float64, long float64, token string) (string, string, error) {
	response, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&exclude=current,minutely,hourly,alerts&units=metric&appid=%s",
		lat,
		long,
		token))
	if err != nil {
		log.Errorf("error response weather api %s", err.Error())
		return "", "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("error parsing weather api answer %s", err.Error())
		return "", "", err
	}

	var resp Weather

	jsonErr := json.Unmarshal(body, &resp)
	if jsonErr != nil {
		log.Errorf("error unmarshaling weather answer: %s", jsonErr)
		return "", "", jsonErr
	}

	return string(fmt.Sprintf("Weather forecast for the day:\nTemperature: min %dC, max %dC, Day %dC, Night %dC\nFeels like: Day %dC, Night %dC\nHumidity: %d%%\nCloudiness: %d%%\nWind speed: %.2fm/s\nWeather: %s(%s)",
		int(resp.Daily[0].Temperature.Min), int(resp.Daily[0].Temperature.Max), int(resp.Daily[0].Temperature.Day), int(resp.Daily[0].Temperature.Night),
		int(resp.Daily[0].FeelsLike.Day), int(resp.Daily[0].FeelsLike.Night), resp.Daily[0].Humidity, resp.Daily[0].Clouds, resp.Daily[0].WindSpeed,
		resp.Daily[0].Weather[0].Main, resp.Daily[0].Weather[0].Description)), resp.TimeZone, nil
}

func GetUserTimezone(message *tgbotapi.Message, token string) (string, error) {
	_, timeZone, err := RequestWeatherForecast(message.Location.Latitude, message.Location.Longitude, token)
	if err != nil {
		log.Errorf("error getting user timezone: %s", err)
		return "", err
	}
	return timeZone, nil
}
