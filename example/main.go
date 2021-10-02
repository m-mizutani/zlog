package main

import "github.com/m-mizutani/zlog"

func main() {
	logger := zlog.New()
	logger.Info("test")
}
