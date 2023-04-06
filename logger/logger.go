package logger

import (
	"subscription-bot/config"

	"go.uber.org/zap"
)

const (
	developmentLevel = "DEBUG"
	productionLevel  = "PRODUCTION"
)

func NewBotLogger(cfg *config.BotConfig) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var sugar = logger.Sugar()
	var err error

	switch cfg.LogLevel {
	case productionLevel:
		logger, err = zap.NewProduction()
		if err != nil {
			sugar.With(err).Errorf("unable to build logger:\n'%v'", err)
			return nil, err
		}
	case developmentLevel:
		logger, err = zap.NewDevelopment()
		if err != nil {
			sugar.With(err).Errorf("unable to build logger:\n'%v'", err)
			return nil, err
		}
	default:
		logger, err = zap.NewDevelopment()
		if err != nil {
			sugar.With(err).Errorf("unable to build logger:\n'%v'", err)
			return nil, err
		}
	}

	sugared := logger.Sugar()

	return sugared, nil
}
