package zlog_test

import "github.com/m-mizutani/zlog"

func ExampleLogLevel() {
	logger := newExampleLogger()

	logger.Level = zlog.LevelInfo
	logger.Debug("debugging")
	logger.Info("information")
	// Output: [info] information
}

func ExampleLogger_SetLogLevel() {
	logger := newExampleLogger()

	logger.SetLogLevel("INFO")

	logger.Debug("debugging")
	logger.Info("information")
	// Output: [info] information
}
