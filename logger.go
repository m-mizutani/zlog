package zlog

import (
	"fmt"

	"github.com/m-mizutani/goerr"
)

type Logger struct {
	Level   LogLevel
	Emitter Emitter
	Filters Filters

	errors []error
	infra  *Infra
}

// New provides default setting zlog logger. Info level, console formatter and stdout.
func New() *Logger {
	return &Logger{
		Level:   LevelInfo,
		Emitter: NewWriter(),

		infra: newInfra(),
	}
}

func (x *Logger) SetLogLevel(level string) error {
	l, err := StrToLogLevel(level)
	if err != nil {
		return err
	}
	x.Level = l
	return nil
}

func (x *Logger) AddFilter(filter Filter) {
	x.Filters = append(x.Filters, filter)
}

func (x *Logger) With(key string, value interface{}) *LogEntity {
	e := x.Log()
	return e.With(key, value)
}

func (x *Logger) Log() *LogEntity {
	return newLogEntity(x)
}

func (x *Logger) Msg(level LogLevel, e *LogEntity, format string, args ...interface{}) {
	if level < x.Level {
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

	if err := x.Emitter.Emit(ev); err != nil {
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
