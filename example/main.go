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
	logger.Writer = os.Stderr
	logger.Info("output to stderr")
}

func changeFormatter() {
	logger := zlog.New()
	logger.Formatter = zlog.NewJsonFormatter()
	logger.Info("output as json format")
}
