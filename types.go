package zlog

import (
	"sort"
	"time"
)

type Log struct {
	Level     LogLevel
	Timestamp time.Time
	Msg       string
	Values    map[string]any
	Error     *Error
}

func (x *Log) OrderedKeys() []string {
	keys := make([]string, 0, len(x.Values))
	for k := range x.Values {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
