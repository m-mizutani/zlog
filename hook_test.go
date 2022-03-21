package zlog_test

import (
	"testing"

	"github.com/m-mizutani/zlog"
	"github.com/stretchr/testify/assert"
)

type testEmitter struct {
	emit func(log *zlog.Log) error
}

func (x *testEmitter) Emit(log *zlog.Log) error {
	return x.emit(log)
}

func TestLogHook(t *testing.T) {
	var calledPre, calledEmit, calledPost int
	e := &testEmitter{}

	logger := zlog.New(zlog.WithEmitter(e),
		zlog.WithPreHook(func(l *zlog.Log) {
			assert.Equal(t, 0, calledPre)
			assert.Equal(t, 0, calledEmit)
			assert.Equal(t, 0, calledPost)
			calledPre++
		}),
		zlog.WithPostHook(func(l *zlog.Log) {
			assert.Equal(t, 1, calledPre)
			assert.Equal(t, 1, calledEmit)
			assert.Equal(t, 0, calledPost)
			calledPost++
		}),
	)

	e.emit = func(log *zlog.Log) error {
		calledEmit++
		assert.Equal(t, 1, calledPre)
		assert.Equal(t, 1, calledEmit)
		assert.Equal(t, 0, calledPost)
		return nil
	}
	logger.Info("test")
	assert.Equal(t, 1, calledPre)
	assert.Equal(t, 1, calledEmit)
	assert.Equal(t, 1, calledPost)
}
