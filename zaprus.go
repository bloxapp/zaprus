package zaprus

import (
	"io"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// NewProxyLogrus returns a muted logrus.Logger that forwards logs to the given zap.Logger
func NewProxyLogrus(logger *zap.Logger) *logrus.Logger {
	lr := logrus.New()
	lr.Out = io.Discard
	lr.Hooks.Add(NewHook(logger))
	return lr
}

// Hook is a logrus.Hook that forwards logs to a zap.Logger
type Hook struct {
	logger *zap.Logger
}

// NewHook returns a Hook for the given zap.Logger
func NewHook(logger *zap.Logger) *Hook {
	return &Hook{logger: logger}
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.PanicLevel:
		hook.logger.Panic(entry.Message, zap.Any("fields", entry.Data))
	case logrus.FatalLevel:
		hook.logger.Fatal(entry.Message, zap.Any("fields", entry.Data))
	case logrus.ErrorLevel:
		hook.logger.Error(entry.Message, zap.Any("fields", entry.Data))
	case logrus.WarnLevel:
		hook.logger.Warn(entry.Message, zap.Any("fields", entry.Data))
	case logrus.InfoLevel:
		hook.logger.Info(entry.Message, zap.Any("fields", entry.Data))
	case logrus.DebugLevel, logrus.TraceLevel:
		hook.logger.Debug(entry.Message, zap.Any("fields", entry.Data))
	default:
		hook.logger.Error(
			"zaprus: unknown level in logrus.Entry",
			zap.String("level", entry.Level.String()),
			zap.String("message", entry.Message),
			zap.Any("fields", entry.Data),
		)
	}
	return nil
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
