package zlog

import (
	"fmt"
	"io"
	"os"

	"github.com/m-mizutani/goerr"
)

type Logger struct {
	Level     LogLevel
	Formatter Formatter
	Writer    io.Writer
	Filters   Filters

	errors []error
	infra  *Infra
}

// New provides default setting zlog logger. Info level, console formatter and stdout.
func New() *Logger {
	return &Logger{
		Level:     LevelInfo,
		Formatter: NewConsoleFormatter(),
		Writer:    os.Stdout,

		infra: newInfra(),
	}
}

func (x *Logger) With(key string, value interface{}) *Entry {
	e := newEntry(x)
	return e.With(key, value)
}

func (x *Logger) Msg(level LogLevel, e *Entry, format string, args ...interface{}) {
	if level < x.Level {
		return // skip
	}

	if e == nil {
		e = &Entry{}
	}
	ev := &Event{
		Level:     level,
		Msg:       fmt.Sprintf(format, args...),
		Timestamp: x.infra.Now(),
		Entry:     *e,
	}

	if err := x.Formatter.Write(ev, x.Writer); err != nil {
		x.errors = append(x.errors, goerr.Wrap(err))
	}
	if _, err := x.Writer.Write([]byte("\n")); err != nil {
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
