package zlog

import (
	"fmt"

	"github.com/m-mizutani/goerr"
)

type Logger struct {
	level   LogLevel
	emitter Emitter
	filters Filters

	errors []error
	infra  *Infra
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
	base := &Logger{
		level:   LevelInfo,
		emitter: NewWriter(),
		infra:   newInfra(),
	}

	for _, opt := range options {
		if err := opt(base); err != nil {
			return nil, err
		}
	}

	return base, nil
}

func (x *Logger) SetLogLevel(level string) error {
	l, err := StrToLogLevel(level)
	if err != nil {
		return err
	}
	x.level = l
	return nil
}

func (x *Logger) ReplaceInfraForTest(infra *Infra) {
	x.infra = infra
}

func (x *Logger) AddFilter(filter Filter) {
	x.filters = append(x.filters, filter)
}

func (x *Logger) With(key string, value interface{}) *LogEntity {
	e := x.Log()
	return e.With(key, value)
}

func (x *Logger) Err(err error) *LogEntity {
	e := x.Log()
	return e.Err(err)
}

func (x *Logger) Log() *LogEntity {
	return newLogEntity(x)
}

func (x *Logger) Msg(level LogLevel, e *LogEntity, format string, args ...interface{}) {
	if level < x.level {
		return // skip
	}

	if e == nil {
		e = &LogEntity{}
	}
	ev := &Event{
		Level:     level,
		Msg:       fmt.Sprintf(format, args...),
		Timestamp: x.infra.Now(),
		LogEntity: *e,
	}

	if err := x.emitter.Emit(ev); err != nil {
		x.errors = append(x.errors, goerr.Wrap(err))
	}
}

func (x *Logger) GetErrors() []error { return x.errors }

func (x *Logger) Trace(format string, args ...interface{}) {
	x.Msg(LevelTrace, nil, format, args...)
}
func (x *Logger) Debug(format string, args ...interface{}) {
	x.Msg(LevelDebug, nil, format, args...)
}
func (x *Logger) Info(format string, args ...interface{}) {
	x.Msg(LevelInfo, nil, format, args...)
}
func (x *Logger) Warn(format string, args ...interface{}) {
	x.Msg(LevelWarn, nil, format, args...)
}
func (x *Logger) Error(format string, args ...interface{}) {
	x.Msg(LevelError, nil, format, args...)
}
