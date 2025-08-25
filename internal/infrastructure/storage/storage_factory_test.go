package storage

import (
	"testing"

	"beers-challenge/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
)

func TestNewRepositoryFactory(t *testing.T) {
	cfg := config.NewConfigProvider()
	factory := NewRepositoryFactory(cfg)
	assert.NotNil(t, factory)
}

func TestCreateBeerRepository(t *testing.T) {
	t.Run("inmemory", func(t *testing.T) {
		cfg := config.NewConfigProvider()
		cfg.GetConfig().Database.Type = "inmemory"
		factory := NewRepositoryFactory(cfg)
		repo, err := factory.CreateBeerRepository()
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("postgres", func(t *testing.T) {
		cfg := config.NewConfigProvider()
		cfg.GetConfig().Database.Type = "postgres"
		factory := NewRepositoryFactory(cfg)
		repo, err := factory.CreateBeerRepository()
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	})

	t.Run("unsupported", func(t *testing.T) {
		cfg := config.NewConfigProvider()
		cfg.GetConfig().Database.Type = "mongodb"
		factory := NewRepositoryFactory(cfg)
		_, err := factory.CreateBeerRepository()
		assert.Error(t, err)
	})

	t.Run("default", func(t *testing.T) {
		cfg := config.NewConfigProvider()
		cfg.GetConfig().Database.Type = "unknown"
		factory := NewRepositoryFactory(cfg)
		repo, err := factory.CreateBeerRepository()
		assert.NoError(t, err)
		assert.NotNil(t, repo) // Should default to inmemory
	})
}

func TestGetSupportedRepositoryTypes(t *testing.T) {
	types := GetSupportedRepositoryTypes()
	assert.Contains(t, types, InMemory)
	assert.Contains(t, types, PostgreSQL)
}

func TestValidateRepositoryType(t *testing.T) {
	assert.NoError(t, ValidateRepositoryType("inmemory"))
	assert.NoError(t, ValidateRepositoryType("postgres"))
	assert.Error(t, ValidateRepositoryType("mongodb"))
}
