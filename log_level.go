package zlog

import (
	"strings"

	"github.com/m-mizutani/goerr"
)

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var strToLevelMap = map[string]LogLevel{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

func (x LogLevel) String() string {
	s, ok := levelToStrMap[x]
	if !ok {
		panic("invalid log level variable")
	}
	return s
}

var levelToStrMap = map[LogLevel]string{}

func init() {
	for k, v := range strToLevelMap {
		levelToStrMap[v] = k
	}
}

func LookupLogLevel(s string) (LogLevel, error) {
	level, ok := strToLevelMap[strings.ToLower(s)]
	if !ok {
		return LevelError, goerr.Wrap(ErrInvalidLogLevel).With("level", s)
	}
	return level, nil
}
