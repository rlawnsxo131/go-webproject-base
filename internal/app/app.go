package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rlawnsxo131/go-webproject-base/internal/config"
	"github.com/rlawnsxo131/go-webproject-base/internal/container"
	"github.com/rlawnsxo131/go-webproject-base/internal/handler"
	"github.com/rlawnsxo131/go-webproject-base/internal/infra/logger"
	"github.com/rlawnsxo131/go-webproject-base/internal/infra/mysql"
	"github.com/rlawnsxo131/go-webproject-base/internal/infra/server"
	"github.com/rlawnsxo131/go-webproject-base/internal/infra/server/middleware"
	"github.com/rlawnsxo131/go-webproject-base/internal/infra/timezone"
)

type App struct {
	cfg *config.Config
	db  *sql.DB
	srv *server.HTTPServer
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}

	logger.SetDefault(cfg.LogLevel)
	timezone.MustSetLocalTimezone(cfg.Timezone)

	database, err := mysql.New()
	if err != nil {
		return nil, fmt.Errorf("mysql initialize error: %w", err)
	}

	mux := http.NewServeMux()
	appContainer := container.New(database)

	healthHandler := handler.NewHealthHandler()
	mux.HandleFunc("GET /health", healthHandler.Get())

	exampleHandler := handler.NewExampleHandler(appContainer)
	mux.HandleFunc("GET /api/example", exampleHandler.Get())
	mux.HandleFunc("GET /api/example/lock", exampleHandler.GetWithLock())
	mux.HandleFunc("POST /api/example", exampleHandler.Create())

	httpHandler := middleware.TraceId(
		middleware.Logger(
			middleware.Recovery(
				middleware.Cors(middleware.CorsConfig{
					AllowedOrigins: []string{"localhost"},
					AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders: []string{"Content-Type", "Authorization"},
					ExposedHeaders: []string{"X-Trace-Id"},
					MaxAge:         3600,
				})(mux),
			),
		),
	)

	srv := server.NewHTTPServer(httpHandler, cfg.Port)

	return &App{
		cfg: cfg,
		db:  database,
		srv: srv,
	}, nil
}

func (a *App) Run() {
	a.srv.Run(func() {
		if err := a.db.Close(); err != nil {
			slog.Error("database close error", "error", err)
		}
	})
}
