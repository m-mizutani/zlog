package zlog_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
)

func newTestLogger(options ...zlog.Option) (*zlog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	emitter := zlog.NewConsoleEmitter(
		zlog.ConsoleNoColor(),
		zlog.ConsoleWriter(buf),
	)

	logger := zlog.New(append(options, zlog.WithEmitter(emitter))...)

	return logger, buf
}

func TestLogger(t *testing.T) {
	t.Run("outout message with values", func(t *testing.T) {
		logger, buf := newTestLogger(zlog.WithClock(func() time.Time {
			return time.Date(2021, 4, 20, 5, 12, 19, 0, time.Local)
		}))
		logger.With("magic", "five").Info("hello %s", "friends")

		msg := buf.String()
		assert.NotContains(t, msg, "2021")
		assert.Contains(t, msg, "05:12:19.000")
		assert.Contains(t, msg, "hello friends")
		assert.Contains(t, msg, "magic")
		assert.Contains(t, msg, "five")
	})

	t.Run("outout message if level is equal or higher than logger level", func(t *testing.T) {
		logger, buf := newTestLogger(zlog.WithLogLevel("warn"))

		logger.Trace("one")
		logger.Debug("two")
		logger.Info("three")
		logger.Warn("four")
		logger.Error("five")

		msg := buf.String()
		assert.NotContains(t, msg, "one")
		assert.NotContains(t, msg, "two")
		assert.NotContains(t, msg, "three")
		assert.Contains(t, msg, "four")
		assert.Contains(t, msg, "five")
	})
}
