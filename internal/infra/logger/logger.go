package logger

import (
	"log/slog"
	"os"
)

func SetDefault(level slog.Level) {
	// multiWriter := io.MultiWriter(os.Stdout, rotateWriter)

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Value.Kind() {
			case slog.KindTime:
				if len(groups) == 0 && a.Key == slog.TimeKey {
					return slog.String(slog.TimeKey, a.Value.Time().Format("2006-01-02 15:04:05.999999999"))
				}
			// µs 형태로 출력되게 만듦
			case slog.KindDuration:
				return slog.String(a.Key, a.Value.Duration().String())
			}

			return a
		},
		// AddSource: true,
	}

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, opts),
		),
	)
}
