package models

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
