package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := GetLogEntry(r.Context())

		defer func() {
			if err := recover(); err != nil {
				// 에러 로그를 남기지 않고 net/http 기본 서버가 조용히 연결을 닫도록 다시 전파
				if err == http.ErrAbortHandler {
					panic(err)
				}

				slog.ErrorContext(r.Context(), "server crash(panic) recovered",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)

				if entry != nil {
					// 이미 상태 코드가 기록되었거나, 본문(Bytes) 데이터가 나갔거나, Flush 등으로 헤더가 전송된 경우
					if entry.ww.Status() != 0 || entry.ww.BytesWritten() > 0 {
						return
					}
				}

				// 101 Switching Protocols로 전환된 웹소켓 연결인 경우 웹소켓 핸들러 내부에서 패닉이 나면 HTTP 응답을 보내면 안됨
				if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"status":500,"data":{"error":{"message":"Internal Server Error"}}}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
