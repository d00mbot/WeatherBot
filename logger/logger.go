package logger

import (
	"subscription-bot/config"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

const (
	developmentLevel = "DEBUG"
	productionLevel  = "PRODUCTION"
)

func NewBotLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	switch cfg.LogLevel {
	case productionLevel:
		logger, err = zap.NewProduction()
		if err != nil {
			log.Errorf("faild to build production logger:\n'%v'", err)
			return nil, err
		}
	case developmentLevel:
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Errorf("faild to build development logger:\n'%v'", err)
			return nil, err
		}
	default:
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Errorf("faild to build development logger:\n'%v'", err)
			return nil, err
		}
	}

	sugared := logger.Sugar()

	return sugared, nil
}
