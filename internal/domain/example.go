package domain

import (
	"errors"
	"time"
)

var (
	ErrExampleInvalidID   = errors.New("invalid ID: must be greater than 0")
	ErrExampleInvalidName = errors.New("invalid Name: must not be empty")
	ErrExampleNotFound    = errors.New("example not found")
)

type Example struct {
	ExampleID int64     `json:"exampleID"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewExample(name string) (*Example, error) {
	if name == "" {
		return nil, ErrExampleInvalidName
	}

	return &Example{
		Name: name,
	}, nil
}

func ReconstructExample(exampleID int64, name string, createdAt time.Time, updatedAt time.Time) (*Example, error) {
	if exampleID <= 0 {
		return nil, ErrExampleInvalidID
	}
	if name == "" {
		return nil, ErrExampleInvalidName
	}

	return &Example{
		ExampleID: exampleID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
