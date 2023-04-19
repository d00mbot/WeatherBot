package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	WeatherToken        string `env:"WEATHER_TOKEN"`
	TelegramToken       string `env:"TELEGRAM_TOKEN"`
	MongoURI            string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017/"`
	LogLevel            string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	ResponsesConfigPath string `env:"RESPONSES_CONFIG_PATH" envDefault:"../config"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Errorf("unable to parse bot config: %s", err)
		return nil, err
	}

	return &cfg, nil
}
