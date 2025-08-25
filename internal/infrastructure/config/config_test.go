package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigProvider(t *testing.T) {
	provider := NewConfigProvider()
	assert.NotNil(t, provider)
	assert.NotNil(t, provider.config)
}

func TestGetString(t *testing.T) {
	os.Setenv("SERVER_HOST", "test_host")
	defer os.Unsetenv("SERVER_HOST")

	provider := NewConfigProvider()
	assert.Equal(t, "test_host", provider.GetString("server.host"))
	assert.Equal(t, "inmemory", provider.GetString("database.type")) // Default
}

func TestGetInt(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	defer os.Unsetenv("SERVER_PORT")

	provider := NewConfigProvider()
	assert.Equal(t, 9090, provider.GetInt("server.port"))
	assert.Equal(t, 5432, provider.GetInt("database.port")) // Default
}

func TestGetDatabaseConnectionString(t *testing.T) {
	provider := NewConfigProvider()
	provider.config.Database.Type = "postgres"
	provider.config.Database.Host = "db_host"
	// ... set other db fields

	connStr := provider.GetDatabaseConnectionString()
	assert.Contains(t, connStr, "host=db_host")
}

func TestIsDevelopment(t *testing.T) {
	os.Setenv("ENVIRONMENT", "development")
	defer os.Unsetenv("ENVIRONMENT")

	provider := NewConfigProvider()
	assert.True(t, provider.IsDevelopment())
	assert.False(t, provider.IsProduction())
}

func TestIsProduction(t *testing.T) {
	os.Setenv("ENVIRONMENT", "production")
	defer os.Unsetenv("ENVIRONMENT")

	provider := NewConfigProvider()
	assert.False(t, provider.IsDevelopment())
	assert.True(t, provider.IsProduction())
}
