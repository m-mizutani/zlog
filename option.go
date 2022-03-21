package zlog

import "time"

type Option func(logger *Logger)

// WithFilters appends filters
func WithFilters(filters ...Filter) Option {
	return func(logger *Logger) {
		logger.filters = append(logger.filters, filters...)
	}
}

// WithLogLevel sets logging level to one of "trace", "debug", "info", "warn", "error" and "fatal". Argument *level* is not case sensitive.
func WithLogLevel(level string) Option {
	return func(logger *Logger) {
		l, err := LookupLogLevel(level)
		if err != nil {
			panic("failed to set log level: " + err.Error())
		}

		logger.level = l
	}
}

// WithEmitter replaces emitter in the logger
func WithEmitter(emitter Emitter) Option {
	return func(logger *Logger) {
		logger.emitter = emitter
	}
}

// WithClock replaces time.Now function in the logger
func WithClock(clock func() time.Time) Option {
	return func(logger *Logger) {
		logger.now = clock
	}
}

// WithErrHook sets hook that is called when emitter has error
func WithErrHook(hook func(error, *Log)) Option {
	return func(logger *Logger) {
		logger.errHooks = append(logger.errHooks, hook)
	}
}

// WithPreHook sets hook that is called before emitting log
func WithPreHook(hook func(*Log)) Option {
	return func(logger *Logger) {
		logger.preHooks = append(logger.preHooks, hook)
	}
}

// WithPostHook sets hook that is called after emitting log
func WithPostHook(hook func(*Log)) Option {
	return func(logger *Logger) {
		logger.postHooks = append(logger.postHooks, hook)
	}
}
