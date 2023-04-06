package config

import (
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
)

type BotConfig struct {
	WeatherToken        string `env:"WEATHER_TOKEN"`
	TelegramToken       string `env:"TELEGRAM_TOKEN"`
	MongoURI            string `env:"MONGO_URI"`
	LogLevel            string `env:"LOG_LEVEL"`
	ResponsesConfigPath string `env:"RESPONSES_CONFIG_PATH"`
}

func LoadConfig() (*BotConfig, error) {
	var cfg BotConfig

	if err := env.Parse(&cfg); err != nil {
		log.Errorf("unable to parse bot config: %s", err)
		return nil, err
	}

	return &cfg, nil
}
