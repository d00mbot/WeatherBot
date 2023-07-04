package models

type Responses struct {
	Start          string `mapstructure:"start"`
	Help           string `mapstructure:"help"`
	UserCreated    string `mapstructure:"user_created"`
	UserUpdated    string `mapstructure:"user_updated"`
	Time           string `mapstructure:"time"`
	WrongTime      string `mapstructure:"wrong_time"`
	TimeUpdated    string `mapstructure:"time_updated"`
	UnknownCommand string `mapstructure:"unknown_command"`
	DefaultMessage string `mapstructure:"default_message"`
}
