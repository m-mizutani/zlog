package zlog

import "time"

type Event struct {
	Level     LogLevel
	Timestamp time.Time
	Msg       string
	Entry
}
