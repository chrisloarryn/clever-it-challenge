package storage

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/infrastructure/storage/inmemory"
	"CleverIT-challenge/internal/infrastructure/storage/postgres"
)

const (
	PROD = "PRODUCTION"
)

func New(environment string) beers.Repository{
	switch environment {
	case PROD:
		return postgres.NewRepository()
	default:
		return inmemory.NewInMemoryRepository()
	}
}
