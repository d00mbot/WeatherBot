package api

type weatherService struct {
	service forecastService
}

type WeatherService interface {
	GetForecast(lat float64, lon float64) (forecast string, timezome string, err error)
}

func NewWeatherService(fs forecastService) WeatherService {
	return &weatherService{service: fs}
}

func (w *weatherService) GetForecast(lat float64, lon float64) (forecast string, timezome string, err error) {
	forecast, timezone, err := w.service.getWeatherForecast(lat, lon)
	if err != nil {
		return "", "", err
	}

	return forecast, timezone, nil
}
