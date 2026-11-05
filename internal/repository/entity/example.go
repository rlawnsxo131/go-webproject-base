package entity

import (
	"time"

	"github.com/rlawnsxo131/go-webproject-base/internal/domain"
)

type Example struct {
	ExampleID int64     `db:"example_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewExampleEntityFromDomain(dom *domain.Example) *Example {
	if dom == nil {
		return nil
	}
	if dom.Name == "" {
		return nil
	}
	return &Example{
		Name: dom.Name,
	}
}

func (e *Example) ToDomain() (*domain.Example, error) {
	domainEntity, err := domain.ReconstructExample(
		e.ExampleID,
		e.Name,
		e.CreatedAt,
		e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return domainEntity, nil
}

// CREATE TABLE example(
//   `example_id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
//   `name` VARCHAR(255) NOT NULL,
//   `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
//   `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//   PRIMARY KEY (`example_id`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
