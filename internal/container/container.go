package container

import (
	"database/sql"
	"sync"

	"github.com/rlawnsxo131/go-webproject-base/internal/repository"
	"github.com/rlawnsxo131/go-webproject-base/internal/service"
)

type Container struct {
	ExampleSvc *service.ExampleService
}

var (
	_onceContainer sync.Once
	container      *Container
)

func New(db *sql.DB) *Container {
	_onceContainer.Do(func() {
		exampleRepo := repository.NewExampleRepository(db)
		exampleSvc := service.NewExampleService(exampleRepo)

		container = &Container{
			ExampleSvc: exampleSvc,
		}
	})
	return container
}
