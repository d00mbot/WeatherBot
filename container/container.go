package container

import (
	"subscription-bot/config"

	"go.uber.org/zap"
)

type container struct {
	config *config.Config
	logger *zap.SugaredLogger
}

type Container interface {
	GetConfig() *config.Config
	GetLogger() *zap.SugaredLogger
}

func NewContainer(cfg *config.Config, log *zap.SugaredLogger) Container {
	return &container{config: cfg, logger: log}
}

func (c *container) GetConfig() *config.Config {
	return c.config
}

func (c *container) GetLogger() *zap.SugaredLogger {
	return c.logger
}
