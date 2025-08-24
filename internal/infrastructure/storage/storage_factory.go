package storage

import (
	"fmt"

	"beers-challenge/internal/core/ports/secondary"
	"beers-challenge/internal/infrastructure/config"
	"beers-challenge/internal/infrastructure/storage/inmemory"
	"beers-challenge/internal/infrastructure/storage/postgres"
)

// RepositoryType represents the type of repository
type RepositoryType string

const (
	InMemory   RepositoryType = "inmemory"
	PostgreSQL RepositoryType = "postgres"
)

// RepositoryFactory creates repositories based on configuration
type RepositoryFactory struct {
	config *config.ConfigProvider
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(config *config.ConfigProvider) *RepositoryFactory {
	return &RepositoryFactory{
		config: config,
	}
}

// CreateBeerRepository creates a beer repository based on the configured database type
func (f *RepositoryFactory) CreateBeerRepository() (secondary.BeerRepository, error) {
	dbType := f.config.GetString("database.type")

	switch RepositoryType(dbType) {
	case PostgreSQL:
		return postgres.NewRepository(f.config)
	case InMemory:
		return inmemory.NewRepository(), nil
	default:
		// Default to in-memory for unknown types
		return inmemory.NewRepository(), nil
	}
}

// GetSupportedRepositoryTypes returns the supported repository types
func GetSupportedRepositoryTypes() []RepositoryType {
	return []RepositoryType{InMemory, PostgreSQL}
}

// ValidateRepositoryType validates if a repository type is supported
func ValidateRepositoryType(repoType string) error {
	for _, supported := range GetSupportedRepositoryTypes() {
		if string(supported) == repoType {
			return nil
		}
	}
	return fmt.Errorf("unsupported repository type: %s", repoType)
}
