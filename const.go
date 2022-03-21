package zlog

import "github.com/m-mizutani/goerr"

const (
	Version       = "v0.0.1"
	FilteredLabel = "[filtered]"
)

var (
	ErrInvalidLogLevel = goerr.New("invalid log level")
)
