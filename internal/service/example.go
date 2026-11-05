package service

import (
	"context"

	"github.com/rlawnsxo131/go-webproject-base/internal/domain"
	"github.com/rlawnsxo131/go-webproject-base/internal/repository"
)

type ExampleService struct {
	exampleRepo *repository.ExampleRepository
}

func NewExampleService(exampleRepo *repository.ExampleRepository) *ExampleService {
	return &ExampleService{
		exampleRepo: exampleRepo,
	}
}

func (s *ExampleService) GetByID(ctx context.Context, exampleID int64) (*domain.Example, error) {
	return s.exampleRepo.FindByID(ctx, exampleID)
}

func (s *ExampleService) GetByIDWithLock(ctx context.Context, id int64) (*domain.Example, error) {
	return s.exampleRepo.FindByID(ctx, id, repository.WithLock())
}

func (s *ExampleService) Create(ctx context.Context, dom *domain.Example) (*domain.Example, error) {
	exampleId, err := s.exampleRepo.Create(ctx, dom)
	if err != nil {
		return nil, err
	}
	return s.exampleRepo.FindByID(ctx, exampleId)
}
