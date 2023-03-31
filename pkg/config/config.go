package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"subscription-bot/pkg/domain/models"
)

type Config struct {
	WeatherToken       string `env:"WEATHER_TOKEN"`
	TelegramToken      string `env:"TELEGRAM_TOKEN"`
	MongoURI           string `env:"MONGO_URI"`
	LogFilePath        string `env:"LOG_FILE_PATH"`
	ResponseConfigPath string `env:"RESPONSE_CONFIG_PATH"`

	Response models.Response
}

func LoadConfig(logger *zap.SugaredLogger) (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		logger.Errorf("unable to load env file: %s", err)
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		logger.Errorf("unable to parse config: %s", err)
		return nil, err
	}

	if err := setUpViper(&cfg, logger); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setUpViper(cfg *Config, logger *zap.SugaredLogger) error {
	viper.AddConfigPath("../../pkg/config") // i have some question here. Will handle this later (env filepath)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("unable to load config file: %s", err)
		return err
	}

	if err := viper.UnmarshalKey("responses", &cfg.Response); err != nil {
		logger.Errorf("faild to unmarshal responses: %s", err)
		return err
	}

	return nil
}
