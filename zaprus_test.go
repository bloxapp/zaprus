package zaprus

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapRecord struct {
	Level   zapcore.Level
	Message string
	Fields  []zapcore.Field
}

type recordingZapCore struct {
	zapcore.Core
	records []zapRecord
}

func (c *recordingZapCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *recordingZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	c.records = append(c.records, zapRecord{
		Level:   entry.Level,
		Message: entry.Message,
		Fields:  fields,
	})
	return c.Core.Write(entry, fields)
}

func TestNewProxy(t *testing.T) {
	z, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	recorder := &recordingZapCore{Core: z.Core()}
	z = zap.New(recorder)

	var expectedRecords []zapRecord
	lr := NewProxy(z)

	lr.WithField("some_key", "some_value").Trace("trace")
	expectedRecords = append(expectedRecords, zapRecord{
		Level:   zapcore.DebugLevel,
		Message: "trace",
		Fields:  []zapcore.Field{zap.Any("fields", logrus.Fields{"some_key": "some_value"})},
	})

	lr.Error("error")
	expectedRecords = append(expectedRecords, zapRecord{
		Level:   zapcore.ErrorLevel,
		Message: "error",
	})

	func() {
		defer func() {
			// Ignore the panic.
			recover()
		}()
		lr.Panic("panic")
	}()
	expectedRecords = append(expectedRecords, zapRecord{
		Level:   zapcore.PanicLevel,
		Message: "panic",
	})

	require.Equal(t, expectedRecords, recorder.records, "expected logs did not match actual logs")
}
