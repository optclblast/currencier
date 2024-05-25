package logger

import (
	"io"
	"log/slog"
	"os"
)

type loggerBuilder struct {
	local     bool
	addSource bool
	lvl       slog.Level
	writers   []io.Writer
}

type NewLoggerOption func(lb *loggerBuilder)

func NewLogger(
	opts ...NewLoggerOption,
) *slog.Logger {
	lb := new(loggerBuilder)

	for _, opt := range opts {
		opt(lb)
	}

	return lb.build()
}

func WithWriter(w io.Writer) NewLoggerOption {
	return func(lb *loggerBuilder) {
		lb.writers = append(lb.writers, w)
	}

}

func WithLevel(l slog.Level) NewLoggerOption {
	return func(lb *loggerBuilder) {
		lb.lvl = l
	}
}

func Local() NewLoggerOption {
	return func(lb *loggerBuilder) {
		lb.local = true

	}
}

func WithSource() NewLoggerOption {
	return func(lb *loggerBuilder) {
		lb.addSource = true

	}
}

func (b *loggerBuilder) build() *slog.Logger {
	w := io.MultiWriter(b.writers...)

	if b.local {
		opts := PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level:     b.lvl,
				AddSource: b.addSource,
			},
		}

		handler := opts.NewPrettyHandler(w)

		return slog.New(handler)
	}

	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     b.lvl,
			AddSource: b.addSource,
		}),
	)
}

func newLogger(lvl slog.Level, w io.Writer) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}),
	)
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func MapLevel(lvl string) slog.Level {
	switch lvl {
	case "dev", "local", "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
