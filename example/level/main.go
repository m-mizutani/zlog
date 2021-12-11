package main

import "github.com/m-mizutani/zlog"

func main() {
	logger := zlog.New(zlog.WithLogLevel("trace"))

	logger.Trace("not")
	logger.Debug("sane")
	logger.Info("five")
	logger.Warn("timeless")
	logger.Error("words")
}
