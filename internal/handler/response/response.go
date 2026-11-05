package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type JSONResponse struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func JSON(ctx context.Context, w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Status: status,
		Data:   data,
	}
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		slog.ErrorContext(ctx, "[response.JSON] json data parsing error", slog.Any("error", err), slog.Any("data", data))
		http.Error(w, `{"status":500,"data":{"error":{"message":"Internal Server Error"}}}`, http.StatusInternalServerError)
	}
}
