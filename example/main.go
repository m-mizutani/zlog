package main

import (
	"os"

	"github.com/m-mizutani/zlog"
)

type myRecord struct {
	Name  string
	EMail string
}

func main() {
	record := myRecord{
		Name:  "mizutani",
		EMail: "mizutani@hey.com",
	}

	logger := zlog.New()
	logger.With("record", record).Info("hello my logger")

	changeWriter()
	changeFormatter()
}

func changeWriter() {
	logger := zlog.New()
	logger.Emitter = zlog.NewWriterWith(zlog.NewConsoleFormatter(), os.Stderr)
	logger.Info("output to stderr")
}

func changeFormatter() {
	logger := zlog.New()
	logger.Emitter = zlog.NewWriterWith(zlog.NewJsonFormatter(), os.Stdout)
	logger.Info("output as json format")
}
