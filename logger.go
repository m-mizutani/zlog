package zlog

import (
	"fmt"
	"reflect"
	"time"
)

type loggerBase struct {
	level     LogLevel
	emitter   Emitter
	filters   Filters
	errHooks  []ErrorHook
	preHooks  []LogHook
	postHooks []LogHook
	now       func() time.Time
}

type Logger struct {
	loggerBase

	err    error
	values map[string]any
}

// New provides default setting zlog logger. Info level, console formatter and stdout.
func New(options ...Option) *Logger {
	logger, err := NewWithError(options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize zlog.Logger: %+v", err))
	}
	return logger
}

func NewWithError(options ...Option) (*Logger, error) {
	logger := &Logger{
		loggerBase: loggerBase{
			level:   LevelInfo,
			emitter: NewConsoleEmitter(),
			now:     time.Now,
		},
		values: make(map[string]any),
	}

	for _, opt := range options {
		opt(logger)
	}

	return logger, nil
}

func (x *Logger) Clone(options ...Option) *Logger {
	newLogger := x.reflect()
	newLogger.filters = x.filters[:]
	newLogger.preHooks = x.preHooks[:]
	newLogger.postHooks = x.postHooks[:]
	newLogger.errHooks = x.errHooks[:]

	for _, opt := range options {
		opt(newLogger)
	}
	return newLogger
}

func (x *Logger) reflect() *Logger {
	newLogger := &Logger{
		loggerBase: x.loggerBase,
		values:     make(map[string]any),
		err:        x.err,
	}
	for k, v := range x.values {
		newLogger.values[k] = v
	}

	return newLogger
}

func (x *Logger) With(key string, value interface{}) *Logger {
	// With sets key-value pair for the log. A previous value is overwritten by same key.
	e := x.reflect()

	if len(x.filters) > 0 && value != nil {
		e.values[key] = newMasking(x.filters).clone(key, reflect.ValueOf(value), "").Interface()
	} else {
		e.values[key] = value
	}

	return e
}

func (x *Logger) Err(err error) *Logger {
	e := x.reflect()
	e.err = err
	return e
}

func (x *Logger) msg(level LogLevel, format string, args ...interface{}) {
	if level < x.level {
		return // skip
	}

	log := &Log{
		Level:     level,
		Msg:       fmt.Sprintf(format, args...),
		Timestamp: x.now(),
		Values:    x.values,
		Error:     newError(x.err),
	}

	for _, hook := range x.preHooks {
		hook(log)
	}
	if err := x.emitter.Emit(log); err != nil {
		for _, hook := range x.errHooks {
			hook(err, log)
		}
	}
	for _, hook := range x.postHooks {
		hook(log)
	}
}

func (x *Logger) Trace(format string, args ...interface{}) {
	x.msg(LevelTrace, format, args...)
}
func (x *Logger) Debug(format string, args ...interface{}) {
	x.msg(LevelDebug, format, args...)
}
func (x *Logger) Info(format string, args ...interface{}) {
	x.msg(LevelInfo, format, args...)
}
func (x *Logger) Warn(format string, args ...interface{}) {
	x.msg(LevelWarn, format, args...)
}
func (x *Logger) Error(format string, args ...interface{}) {
	x.msg(LevelError, format, args...)
}
func (x *Logger) Fatal(format string, args ...interface{}) {
	x.msg(LevelFatal, format, args...)
}
