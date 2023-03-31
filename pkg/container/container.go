package container

import (
	"subscription-bot/pkg/config"

	"go.uber.org/zap"
)

type container struct {
	config *config.Config
	logger *zap.SugaredLogger
}

type BotContainer interface {
	GetConfig() *config.Config
	GetLogger() *zap.SugaredLogger
}

func NewBotContainer(cfg *config.Config, log *zap.SugaredLogger) BotContainer {
	return &container{config: cfg, logger: log}
}

func (c *container) GetConfig() *config.Config {
	return c.config
}

func (c *container) GetLogger() *zap.SugaredLogger {
	return c.logger
}
