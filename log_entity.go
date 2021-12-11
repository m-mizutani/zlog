package zlog

import (
	"reflect"
)

type LogEntity struct {
	logger *Logger
	values map[string]interface{}
	err    *Error
}

func newLogEntity(logger *Logger) *LogEntity {
	return &LogEntity{
		logger: logger,
		values: make(map[string]interface{}),
	}
}

// Values copies key-value map and return it
func (x *LogEntity) Values() map[string]interface{} {
	kv := map[string]interface{}{}
	for key, value := range x.values {
		kv[key] = value
	}
	return kv
}

func (x *LogEntity) Clone() *LogEntity {
	entry := newLogEntity(x.logger)
	for k, v := range x.values {
		entry.values[k] = v
	}
	entry.err = x.err

	return entry
}

// With sets key-value pair for the log. A previous value is overwritten by same key.
func (x *LogEntity) With(key string, value interface{}) *LogEntity {
	e := x.Clone()

	if len(x.logger.filters) > 0 && value != nil {
		e.values[key] = newMasking(x.logger.filters).clone(key, reflect.ValueOf(value), "").Interface()
	} else {
		e.values[key] = value
	}

	return e
}

// Err sets error and extracts stacktrace if available. LogEntity can have only one error
func (x *LogEntity) Err(err error) *LogEntity {
	e := x.Clone()
	e.err = newError(err)
	return e
}

func (x *LogEntity) Trace(format string, args ...interface{}) {
	x.logger.Msg(LevelTrace, x, format, args...)
}
func (x *LogEntity) Debug(format string, args ...interface{}) {
	x.logger.Msg(LevelDebug, x, format, args...)
}
func (x *LogEntity) Info(format string, args ...interface{}) {
	x.logger.Msg(LevelInfo, x, format, args...)
}
func (x *LogEntity) Warn(format string, args ...interface{}) {
	x.logger.Msg(LevelWarn, x, format, args...)
}
func (x *LogEntity) Error(format string, args ...interface{}) {
	x.logger.Msg(LevelError, x, format, args...)
}
