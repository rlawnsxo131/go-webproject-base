package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(handler http.Handler, port string) *HTTPServer {
	errorLogger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)
	errorLogger.SetPrefix("[SVR_ERR] ")

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", "0.0.0.0", port),

		// 응답 데이터 전송 최대 시간
		WriteTimeout: 15 * time.Second,
		// 요청 헤더를 읽는 최대 시간
		ReadHeaderTimeout: 5 * time.Second,
		// 요청 바디를 읽는 최대 시간
		ReadTimeout: 15 * time.Second,
		// 연결(Keep-Alive)을 유지하는 최대 시간
		IdleTimeout: 60 * time.Second,

		Handler: handler,

		// 서버 인프라 자체의 내부 결함이나 네트워크 레벨의 문제로 발생하는 에러 기록
		// 내부 HTTP 서버(네트워크 연결 처리, 프로토콜 위반, 헤더 파싱 오류 등)가 구동 중에 발생시키는 자체 에러 기록
		ErrorLog: errorLogger,
	}
	srv.SetKeepAlivesEnabled(true)

	return &HTTPServer{
		srv,
	}
}

func (s *HTTPServer) Run(doShutdown func()) {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxCancel()

		if doShutdown != nil {
			doShutdown()
		}

		if err := s.Shutdown(shutdownCtx); err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				slog.Info("graceful shutdown timed out. forcing exit.")
			case errors.Is(err, context.Canceled):
				slog.Info("graceful shutdown context was canceled.")
			default:
				slog.Error("HTTP Server Shutdown Failed", slog.Any("error", err))
			}
		}

		serverStopCtx()
	}()

	slog.Info("starting HTTP server", slog.String("addr", s.Addr))
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("failed to start server", slog.Any("error", err))
		panic(err)
	}

	<-serverCtx.Done()
	slog.Info("shutting down HTTP server")
}
