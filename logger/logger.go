package logger

import (
	"os"
	"subscription-bot/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	sugared := logger.Sugar()

	if err != nil {
		sugared.With(err).Errorf("unable to build logger: %s", err)
		return nil, err
	}

	return sugared, nil
}

func LoadLogger(cfg *config.Config) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileWriter, _, err := zap.Open(cfg.LogFilePath)
	if err != nil {
		return nil, err
	}

	fileWriteSyncer := zapcore.AddSync(fileWriter)

	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	writeSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)

	config.OutputPaths = []string{cfg.LogFilePath}
	config.ErrorOutputPaths = []string{cfg.LogFilePath}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		writeSyncer,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger.Sugar(), nil
}
