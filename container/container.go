package container

import (
	"subscription-bot/config"

	"go.uber.org/zap"
)

type container struct {
	config *config.BotConfig
	logger *zap.SugaredLogger
}

type BotContainer interface {
	GetConfig() *config.BotConfig
	GetLogger() *zap.SugaredLogger
}

func NewBotContainer(cfg *config.BotConfig, log *zap.SugaredLogger) BotContainer {
	return &container{config: cfg, logger: log}
}

func (c *container) GetConfig() *config.BotConfig {
	return c.config
}

func (c *container) GetLogger() *zap.SugaredLogger {
	return c.logger
}
