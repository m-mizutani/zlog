package zlog

import "time"

type Log struct {
	Level     LogLevel
	Timestamp time.Time
	Msg       string
	Values    map[string]interface{}
	Error     *Error
}
