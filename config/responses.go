package config

import (
	"subscription-bot/internal/models"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadResponses(cfg *BotConfig, logger *zap.SugaredLogger) (*models.Responses, error) {
	var rsp models.Responses

	if err := setUpViper(cfg, logger); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("responses", &rsp); err != nil {
		logger.Errorf("faild to unmarshal responses: %s", err)
		return nil, err
	}

	return &rsp, nil
}

func setUpViper(cfg *BotConfig, logger *zap.SugaredLogger) error {
	viper.AddConfigPath(cfg.ResponsesConfigPath)
	viper.SetConfigName("responses")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("unable to load responses config file: %s", err)
		return err
	}

	return nil
}
