package x

import (
	"context"
	"errors"
	"time"
)

func Retry[T any](ctx context.Context, attempts int, sleep time.Duration, fn func() (T, error)) (T, error) {
	res, err := fn()
	if err != nil {
		// 만약 에러의 원인이 컨텍스트 만료/취소라면, 더 이상 재시도하지 않고 즉시 원본 에러를 반환
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			var zero T
			return zero, err
		}

		if attempts--; attempts > 0 {
			// time.Sleep(sleep) 대신 select를 사용하여 대기하는 도중에도
			// 상위 웹 서버나 클라이언트가 요청을 취소하면 즉시 깨어나서 종료되도록 방어
			select {
			case <-ctx.Done():
				// 대기 중 컨텍스트가 끝나면 즉시 컨텍스트 에러 반환
				var zero T
				return zero, ctx.Err()
			case <-time.After(sleep):
				// 지정된 지수 백오프 시간만큼 무사히 대기했다면 아래 블록으로 넘어감(재시도)
			}

			return Retry(ctx, attempts, 2*sleep, fn)
		}

		var zero T
		return zero, err
	}

	return res, nil
}
