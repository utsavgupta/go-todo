package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type zeroLogger struct {
	instance zerolog.Logger
}

func NewZeroLogger() Logger {

	instance := zerolog.New(os.Stdout)
	return &zeroLogger{instance}
}

func (logger *zeroLogger) Info(ctx context.Context, message string) {

	logger.instance.Info().Msg(message)
}

func (logger *zeroLogger) Debug(ctx context.Context, message string) {

	logger.instance.Debug().Msg(message)
}

func (logger *zeroLogger) Warn(ctx context.Context, message string) {

	logger.instance.Warn().Msg(message)
}

func (logger *zeroLogger) Error(ctx context.Context, message string) {

	logger.instance.Error().Msg(message)
}
