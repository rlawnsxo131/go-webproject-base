package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/rlawnsxo131/go-webproject-base/internal/infra/server/middleware/response"
)

var _loggerCtxKey = newContextKey("logger")

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := response.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		entry := &logEntry{r: r, ww: ww}

		defer func() {
			entry.write(start)
		}()

		next.ServeHTTP(ww, withLogEntry(r, entry))
	})
}

type logEntry struct {
	r     *http.Request
	ww    response.WrapResponseWriter
	extra []slog.Attr
}

func (e *logEntry) write(start time.Time) {
	// infra/logger/logger.go 에서 slog.SetDefault() 로 설정한 로그 레벨이 info 보다 낮으면, 로그를 남기지 않음
	if !slog.Default().Enabled(e.r.Context(), slog.LevelInfo) {
		return
	}

	attrs := make([]slog.Attr, 0, 8+len(e.extra))

	attrs = append(attrs,
		slog.String("trace_id", GetTraceId(e.r.Context())),
		slog.Int("status", e.ww.Status()),
		slog.String("method", e.r.Method),
		slog.String("path", e.r.URL.Path),
		slog.String("remote_addr", e.r.RemoteAddr),
		slog.String("user_agent", e.r.UserAgent()),
		slog.Int("byte", e.ww.BytesWritten()),
		slog.Duration("latency", time.Since(start)),
	)
	attrs = append(attrs, e.extra...)

	slog.LogAttrs(e.r.Context(), slog.LevelInfo, "http_request", attrs...)
}

// go routine 내부에서 Add 를 호출하지 않는다는 전제하에 lock 을 잡지 않는다
func (e *logEntry) Add(attrs ...slog.Attr) {
	// infra/logger/logger.go 에서 slog.SetDefault() 로 설정한 로그 레벨이 info 보다 낮으면, 로그를 남기지 않음
	if !slog.Default().Enabled(e.r.Context(), slog.LevelInfo) {
		return
	}

	// 기존 슬라이스의 여유 공간이 부족하면, 필요한 크기만큼 단 한 번만 확장
	needed := len(attrs)
	available := cap(e.extra) - len(e.extra)

	if available < needed {
		newExtra := make([]slog.Attr, len(e.extra), len(e.extra)+needed)
		copy(newExtra, e.extra)
		e.extra = newExtra
	}

	// 루프를 돌며 개별 append 하지 않고, 슬라이스를 한 번에 결합하는게 속도가 가장 빠름
	e.extra = append(e.extra, attrs...)
}

func withLogEntry(r *http.Request, entry *logEntry) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), _loggerCtxKey, entry))
	return r
}

func GetLogEntry(ctx context.Context) *logEntry {
	if entry, _ := ctx.Value(_loggerCtxKey).(*logEntry); entry != nil {
		return entry
	}
	return nil
}
