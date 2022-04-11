package zlog_test

import (
	"testing"
	"time"

	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
)

type delayWriter struct {
	output []byte
}

func (x *delayWriter) Write(data []byte) (int, error) {
	time.Sleep(time.Second)
	x.output = append(x.output, data...)
	return len(data), nil
}

func TestAsync(t *testing.T) {
	buf := &delayWriter{}
	logger := zlog.New(
		zlog.WithAsync(128),
		zlog.WithEmitter(
			zlog.NewJsonEmitter(
				zlog.JsonWriter(buf),
			),
		),
	)

	logger.Info("blue")
	assert.NotContains(t, string(buf.output), "blue")
	logger.Flush()
	assert.Contains(t, string(buf.output), "blue")
}
