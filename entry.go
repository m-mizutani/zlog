package zlog

type Entry struct {
	logger *Logger
	values map[string]interface{}
}

func newEntry(logger *Logger) *Entry {
	return &Entry{
		logger: logger,
		values: make(map[string]interface{}),
	}
}

// Values copies key-value map and return it
func (x *Entry) Values() map[string]interface{} {
	kv := map[string]interface{}{}
	for key, value := range x.values {
		kv[key] = value
	}
	return kv
}

// With sets key-value pair for the log. A previous value is overwritten by same key.
func (x *Entry) With(key string, value interface{}) *Entry {
	x.values[key] = value
	return x
}

func (x *Entry) Trace(format string, args ...interface{}) {
	x.logger.Msg(LevelTrace, x, format, args...)
}
func (x *Entry) Debug(format string, args ...interface{}) {
	x.logger.Msg(LevelDebug, x, format, args...)
}
func (x *Entry) Info(format string, args ...interface{}) {
	x.logger.Msg(LevelInfo, x, format, args...)
}
func (x *Entry) Warn(format string, args ...interface{}) {
	x.logger.Msg(LevelWarn, x, format, args...)
}
func (x *Entry) Error(format string, args ...interface{}) {
	x.logger.Msg(LevelError, x, format, args...)
}
