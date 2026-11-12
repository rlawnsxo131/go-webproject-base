package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/rlawnsxo131/go-webproject-base/internal/container"
	"github.com/rlawnsxo131/go-webproject-base/internal/domain"
	"github.com/rlawnsxo131/go-webproject-base/internal/handler/response"
)

type exampleHandler struct {
	container *container.Container
}

func NewExampleHandler(container *container.Container) *exampleHandler {
	return &exampleHandler{
		container: container,
	}
}

func (h *exampleHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		example, err := h.container.ExampleSvc.GetByID(r.Context(), 1)

		if err != nil {
			if errors.Is(err, domain.ErrExampleNotFound) {
				response.JSON(r.Context(), w, http.StatusNotFound, nil)
				return
			}
			slog.ErrorContext(r.Context(), "[exampleHandler.Get] example get error", slog.Any("error", err))
			response.JSON(r.Context(), w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(r.Context(), w, http.StatusOK, example)
	}
}

func (h *exampleHandler) GetWithLock() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		example, err := h.container.ExampleSvc.GetByIDWithLock(r.Context(), 1)

		if err != nil {
			if errors.Is(err, domain.ErrExampleNotFound) {
				response.JSON(r.Context(), w, http.StatusNotFound, nil)
				return
			}
			slog.ErrorContext(r.Context(), "[exampleHandler.GetWithLock] example get error", slog.Any("error", err))
			response.JSON(r.Context(), w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(r.Context(), w, http.StatusOK, example)
	}
}

//	curl -X POST http://localhost:8080/api/example \
//	  -H "Content-Type: application/json" \
//	  -d '{"name": "테스트이름"}'
func (h *exampleHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.ErrorContext(r.Context(), "[exampleHandler.Create] json decode error", slog.Any("error", err))
			response.JSON(r.Context(), w, http.StatusBadRequest, nil)
			return
		}

		newExample, err := domain.NewExample(req.Name)
		if err != nil {
			if errors.Is(err, domain.ErrExampleInvalidName) {
				response.JSON(r.Context(), w, http.StatusBadRequest, nil)
				return
			}
			slog.ErrorContext(r.Context(), "[exampleHandler.Create] example create error", slog.Any("error", err))
			response.JSON(r.Context(), w, http.StatusInternalServerError, nil)
			return
		}

		example, err := h.container.ExampleSvc.Create(r.Context(), newExample)
		if err != nil {
			slog.ErrorContext(r.Context(), "[exampleHandler.Create] example create error", slog.Any("error", err))
			response.JSON(r.Context(), w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(r.Context(), w, http.StatusCreated, example)
	}
}
