package zlog

import "sync"

type async struct {
	ch chan *asyncMsg
	wg sync.WaitGroup
}

type asyncMsg struct {
	logger *Logger
	log    *Log
}

func newAsync(queueSize int) *async {
	x := &async{
		ch: make(chan *asyncMsg, queueSize),
	}
	x.wg.Add(1)
	go func() {
		defer x.wg.Done()
		for msg := range x.ch {
			msg.logger.emit(msg.log)
		}
	}()

	return x
}

func (x *async) emit(logger *Logger, log *Log) {
	x.ch <- &asyncMsg{
		logger: logger,
		log:    log,
	}
}

func (x *async) flush() {
	close(x.ch)
	x.wg.Wait()
}
