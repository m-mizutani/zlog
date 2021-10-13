package zlog

import "reflect"

type LogEntity struct {
	logger *Logger
	values map[string]interface{}
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

// With sets key-value pair for the log. A previous value is overwritten by same key.
func (x *LogEntity) With(key string, value interface{}) *LogEntity {
	e := newLogEntity(x.logger)
	for k, v := range x.values {
		e.values[k] = v
	}

	if len(x.logger.Filters) > 0 && value != nil {
		e.values[key] = newCensor(x.logger.Filters).clone(reflect.ValueOf(value), "").Interface()
	} else {
		e.values[key] = value
	}

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
