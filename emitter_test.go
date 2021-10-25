package zlog_test

import (
	"fmt"

	"github.com/m-mizutani/zlog"
)

type myEmitter struct {
	seq int
}

func (x *myEmitter) Emit(ev *zlog.Event) error {
	x.seq++
	prefix := []string{"＼(^o^)／", "(´・ω・｀)", "(・∀・)"}
	fmt.Println(prefix[x.seq%3], ev.Msg)
	return nil
}

func ExampleEmitter() {
	logger := zlog.New()
	logger.Emitter = &myEmitter{}

	logger.Info("waiwai")
	logger.Info("heyhey")
	// Output:
	// ＼(^o^)／ waiwai
	// (´・ω・｀) heyhey
}
