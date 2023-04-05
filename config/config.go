package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"subscription-bot/internal/models"
)

type Config struct {
	WeatherToken        string `env:"WEATHER_TOKEN"`
	TelegramToken       string `env:"TELEGRAM_TOKEN"`
	MongoURI            string `env:"MONGO_URI"`
	LogFilePath         string `env:"LOG_FILE_PATH"`
	ResponsesConfigPath string `env:"RESPONSES_CONFIG_PATH"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		log.Errorf("unable to load env file: %s", err)
		return nil, err
	}

	if err := env.Parse(&cfg); err != nil {
		log.Errorf("unable to parse config: %s", err)
		return nil, err
	}

	return &cfg, nil
}

func LoadResponses(cfg *Config) (*models.Responses, error) {
	var rsp models.Responses

	if err := setUpViper(cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("responses", &rsp); err != nil {
		log.Errorf("faild to unmarshal responses: %s", err)
		return nil, err
	}

	return &rsp, nil
}

func setUpViper(cfg *Config) error {
	viper.AddConfigPath(cfg.ResponsesConfigPath) // i have some question here. Will handle this later (env filepath)
	viper.SetConfigName("responses")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("unable to load config file: %s", err)
		return err
	}

	return nil
}
