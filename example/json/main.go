package main

import (
	"github.com/m-mizutani/zlog"
	"github.com/m-mizutani/zlog/filter"
)

type request struct {
	Name     string
	Password string `zlog:"secure"`
}

func main() {
	logger := zlog.New(
		zlog.WithEmitter(zlog.NewJsonEmitter(
			zlog.JsonPrettyPrint(),
		)),
		zlog.WithFilters(filter.Tag("secure")),
	)

	req := &request{
		Name:     "mizutani",
		Password: "abc123",
	}

	logger.With("req", req).Info("send request")
}
