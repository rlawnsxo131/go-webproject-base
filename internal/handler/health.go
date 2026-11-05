package handler

import "net/http"

type healthHandler struct{}

func NewHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
