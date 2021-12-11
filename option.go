package zlog

type Option func(logger *Logger) error

func WithFilters(filters ...Filter) Option {
	return func(logger *Logger) error {
		logger.filters = append(logger.filters, filters...)
		return nil
	}
}

func WithLogLevel(level string) Option {
	return func(logger *Logger) error {
		l, err := StrToLogLevel(level)
		if err != nil {
			return err
		}

		logger.level = l
		return nil
	}
}

func WithEmitter(emitter Emitter) Option {
	return func(logger *Logger) error {
		logger.emitter = emitter
		return nil
	}
}
