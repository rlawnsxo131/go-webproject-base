package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func Cors(config CorsConfig) func(http.Handler) http.Handler {
	allowedMethods := strings.Join(config.AllowedMethods, ",")
	allowedHeaders := strings.Join(config.AllowedHeaders, ",")
	exposedHeaders := strings.Join(config.ExposedHeaders, ",")
	maxAge := "0"
	if config.MaxAge > 0 {
		maxAge = strconv.Itoa(config.MaxAge)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowOrigin := ""
			preflight := r.Method == http.MethodOptions

			if origin == "" {
				if !preflight {
					// browser 의경우 Sec-Fetch-Mode 검사를 할까했으나, 그냥 csrf 를 적용하는것으로 결정
					// server to server
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			w.Header().Add("Vary", "Origin")

			for _, allowedOrigin := range config.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowOrigin = allowedOrigin
					break
				}
			}

			// origin not allowed
			if allowOrigin == "" {
				if !preflight {
					// origin 이 아예 없고 preflight 가 아닌 경우는 server to server 로 판단하고, 그냥 요청을 처리
					// 중요 api 는 각 handler 에 접근전에 도메인 특성에 따라 검증해서 접근을 막는다
					next.ServeHTTP(w, r)
					return
				}
				// 실제 차단은 브라우저에 위임
				w.WriteHeader(http.StatusNoContent)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			if config.AllowCredentials {
				if allowOrigin != "*" {
					// preflight 와 본요청 모두 필요함
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// simple request
			if !preflight {
				if exposedHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
				}
				// 세션/쿠키를 사용하는 데이터 응답인 경우 캐시컨트롤 추가
				// w.Header().Set("Cache-Control", "no-store")
				next.ServeHTTP(w, r)
				return
			}

			// prefight
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)

			if allowedHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			} else {
				if headers := r.Header.Get("Access-Control-Request-Headers"); headers != "" {
					w.Header().Set("Access-Control-Allow-Headers", headers)
				}
			}

			if maxAge != "0" {
				w.Header().Set("Access-Control-Max-Age", maxAge)
			}

			w.WriteHeader(http.StatusNoContent)
		})
	}
}
