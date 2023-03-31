package models

type Response struct {
	Start          string `mapstructure:"start"`
	Help           string `mapstructure:"help"`
	Location       string `mapstructure:"location"`
	LocationUpdate string `mapstructure:"location_update"`
	Time           string `mapstructure:"time"`
	WrongTime      string `mapstructure:"wrong_time"`
	SuccessfulTime string `mapstructure:"successful_time"`
	UnknownCommand string `mapstructure:"unknown_command"`
	DefaultMessage string `mapstructure:"default_message"`
}
