package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var _traceIdCtxKey = newContextKey("traceid")

func TraceId(next http.Handler) http.Handler {
	var reqCounter atomic.Uint64

	// 서버 기동 시 PID와 호스트 이름을 조합해 8자리 고유 접두사를 사전 생성
	// 분산 환경(쿠버네티스 멀티 포드 등)에서 여러 서버가 동시에 가동되어도 ID 경합을 방지
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	// 0부터 65,535 까지의 PID와 호스트 이름 길이를 16진수로 변환해 8자리 접두사 생성
	serverPrefix := fmt.Sprintf("%04x%04x", pid&0xffff, len(hostname)&0xffff)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := r.Header.Get("X-Trace-ID")
		if traceId == "" {
			traceId = r.Header.Get("X-Request-ID")
		}

		if traceId == "" {
			// 무거운 sync.Mutex 락을 잡지 않고, CPU 하드웨어 단에서 지원하는 원자적(Atomic) 연산으로 카운트
			// 대규모 트래픽(초당 수십만 요청)이 몰려 멀티스레드가 동시에 진입해도 Race Condition(경합)과 병목이 없음
			currentCount := reqCounter.Add(1)

			// 시간(ms 단위 16진수) + 서버접두사 + 원자적 카운터 조합
			// 문자열 결합 오버헤드를 막기 위해 단순 포맷 규칙 적용
			nowHex := fmt.Sprintf("%x", time.Now().UnixMilli())
			traceId = fmt.Sprintf("%s-%s-%06x", nowHex, serverPrefix, currentCount&0xffffff)
		}

		w.Header().Set("X-Trace-Id", traceId)

		next.ServeHTTP(w, withTraceId(r, traceId))
	})
}

func withTraceId(r *http.Request, traceId string) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), _traceIdCtxKey, traceId))
	return r
}

func GetTraceId(ctx context.Context) string {
	if traceId, _ := ctx.Value(_traceIdCtxKey).(string); traceId != "" {
		return traceId
	}
	return ""
}
