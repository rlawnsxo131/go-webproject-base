package main

import (
	"log/slog"

	"github.com/rlawnsxo131/go-webproject-base/internal/app"

	// ecs 위에서 돌릴땐: https://github.com/rdforte/gomaxecs
	_ "go.uber.org/automaxprocs"
)

func main() {
	app, err := app.New()

	if err != nil {
		slog.Error("app initialize error", "error", err)
		panic(err)
	}

	app.Run()
}

// docker run -d --cpus=2 -p 8080:8080
