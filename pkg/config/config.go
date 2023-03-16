package config

import (
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Responses struct {
	Start          string `mapstructure:"start"`
	Help           string `mapstructure:"help"`
	Location       string `mapstructure:"location"`
	Time           string `mapstructure:"time"`
	WrongTime      string `mapstructure:"wrong_time"`
	SuccessfulTime string `mapstructure:"successful_time"`
	UnknownCommand string `mapstructure:"unknown_command"`
	DefaultMessage string `mapstructure:"default_message"`
}

type Config struct {
	WeatherToken  string
	TelegramToken string
	MongoURI      string

	Responses Responses
}

func Init() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("faild to load responses config file: %s", err)
		return err
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Errorf("error unmarshaling config into a struct: %s", err)
		return err
	}

	if err := viper.UnmarshalKey("responses", &cfg); err != nil {
		log.Errorf("error unmarshaling key into a struct: %s", err)
		return err
	}

	return nil
}

func fromEnv(cfg *Config) error {
	if err := viper.BindEnv("telegram_token"); err != nil {
		log.Errorf("error binding a Viper key to a ENV variable: %s", err)
		return err
	}
	cfg.TelegramToken = viper.GetString("telegram_token")

	if err := viper.BindEnv("weather_token"); err != nil {
		log.Errorf("error binding a Viper key to a ENV variable: %s", err)
		return err
	}
	cfg.WeatherToken = viper.GetString("weather_token")

	if err := viper.BindEnv("mongo_uri"); err != nil {
		log.Errorf("error binding a Viper key to a ENV variable: %s", err)
		return err
	}
	cfg.MongoURI = viper.GetString("mongo_uri")

	return nil
}
