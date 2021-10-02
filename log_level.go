package zlog

import "strings"

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

var strToLevelMap = map[string]LogLevel{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

func (x LogLevel) String() string {
	s, ok := levelToStrMap[x]
	if !ok {
		panic("invalid log level variable")
	}
	return s
}

var levelToStrMap = map[LogLevel]string{}

func initLogLevelMap() {
	for k, v := range strToLevelMap {
		levelToStrMap[v] = k
	}
}

func StrToLogLevel(s string) (LogLevel, error) {
	level, ok := strToLevelMap[strings.ToLower(s)]
	if !ok {
		return LevelError, nil
	}
	return level, nil
}
