package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type Responses struct {
	Start          string `mapstructure:"start"`
	Time           string `mapstructure:"time"`
	WrongTime      string `mapstructure:"wrong_time"`
	SuccessfulTime string `mapstructure:"successful_time"`
	UnknownCommand string `mapstructure:"unknown_command"`
	Location       string `mapstructure:"location"`
	OtherMessage   string `mapstructure:"other_message"`
}

type Config struct {
	WeatherToken  string `env:"WEATHER_TOKEN"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	MongoURI      string `env:"MONGO_URI"`

	Responses Responses
}

func InitResponses() (*Config, error) {
	if err := setUpViper(); err != nil {
		return nil, err
	}

	var rsp Config

	if err := unmarshal(&rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
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

func unmarshal(rsp *Config) error {
	if err := viper.Unmarshal(&rsp); err != nil {
		log.Errorf("error unmarshaling config into a struct: %s", err)
		return err
	}

	if err := viper.UnmarshalKey("responses", &rsp); err != nil {
		log.Errorf("error unmarshaling key into a struct: %s", err)
		return err
	}

	return nil
}

func LoadENV(fileName string) (*Config, error) {
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("unable to load .env file: %s", err)
		return nil, err
	}

	cfg := Config{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse .env variables %s: ", err)
		return nil, err
	}

	return &cfg, nil
}
