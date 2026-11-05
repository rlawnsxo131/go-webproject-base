package logger

import (
	"log/slog"
	"os"
	"time"
)

func SetDefault(level slog.Level) {
	// multiWriter := io.MultiWriter(os.Stdout, rotateWriter)

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, t.Format("2006-01-02 15:04:05.999999999"))
				}
			}

			// µs 형태로 출력되게 만듦
			if a.Value.Kind() == slog.KindDuration {
				d := a.Value.Duration()
				return slog.String(a.Key, d.String())
			}

			return a
		},
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, opts),
		),
	)
}
