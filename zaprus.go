package zaprus

import (
	"io"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// NewProxy returns a muted logrus.Logger that forwards logs to the given zap.Logger
func NewProxy(to *zap.Logger) (from *logrus.Logger) {
	from = logrus.New()
	from.Level = logrus.TraceLevel
	Mute(from)
	from.Hooks.Add(NewHook(to))
	return from
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
	// Map logrus.Level to zapcore.Level
	var fn func(msg string, fields ...zap.Field)
	switch entry.Level {
	case logrus.PanicLevel:
		fn = hook.logger.Panic
	case logrus.FatalLevel:
		fn = hook.logger.Fatal
	case logrus.ErrorLevel:
		fn = hook.logger.Error
	case logrus.WarnLevel:
		fn = hook.logger.Warn
	case logrus.InfoLevel:
		fn = hook.logger.Info
	case logrus.DebugLevel, logrus.TraceLevel:
		fn = hook.logger.Debug
	default:
		hook.logger.Error(
			"zaprus: unknown level in logrus.Entry",
			zap.String("level", entry.Level.String()),
			zap.String("message", entry.Message),
			zap.Any("fields", entry.Data),
		)
		return nil
	}

	// Log with zap.Logger
	var fields []zap.Field
	if len(entry.Data) > 0 {
		fields = []zap.Field{zap.Any("fields", entry.Data)}
	}
	fn(entry.Message, fields...)

	return nil
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Mute disables the output of the given logrus.Logger
func Mute(lr *logrus.Logger) {
	lr.Formatter = NopFormatter{}
	lr.Out = io.Discard
}

// NopFormatter is a logrus.Formatter that does nothing,
// useful when disabling logrus output.
type NopFormatter struct{}

func (f NopFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return nil, nil
}
