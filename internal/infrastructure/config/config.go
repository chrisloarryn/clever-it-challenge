package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Currency CurrencyConfig `json:"currency"`
	Logger   LoggerConfig   `json:"logger"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"ssl_mode"`
}

// CurrencyConfig holds currency service configuration
type CurrencyConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Timeout int    `json:"timeout"`
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// ConfigProvider implements the secondary.ConfigProvider interface
type ConfigProvider struct {
	config *Config
}

// NewConfigProvider creates a new config provider
func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{
		config: loadConfig(),
	}
}

// GetString returns a string configuration value
func (c *ConfigProvider) GetString(key string) string {
	switch key {
	case "server.host":
		return c.config.Server.Host
	case "database.type":
		return c.config.Database.Type
	case "database.host":
		return c.config.Database.Host
	case "database.name":
		return c.config.Database.Name
	case "database.user":
		return c.config.Database.User
	case "database.password":
		return c.config.Database.Password
	case "database.ssl_mode":
		return c.config.Database.SSLMode
	case "currency.api_key":
		return c.config.Currency.APIKey
	case "currency.base_url":
		return c.config.Currency.BaseURL
	case "logger.level":
		return c.config.Logger.Level
	case "logger.format":
		return c.config.Logger.Format
	default:
		return ""
	}
}

// GetInt returns an integer configuration value
func (c *ConfigProvider) GetInt(key string) int {
	switch key {
	case "server.port":
		return c.config.Server.Port
	case "database.port":
		return c.config.Database.Port
	case "currency.timeout":
		return c.config.Currency.Timeout
	default:
		return 0
	}
}

// GetFloat64 returns a float64 configuration value
func (c *ConfigProvider) GetFloat64(key string) float64 {
	// No float64 configs currently, but method required by interface
	return 0.0
}

// GetBool returns a boolean configuration value
func (c *ConfigProvider) GetBool(key string) bool {
	// No boolean configs currently, but method required by interface
	return false
}

// GetConfig returns the full configuration
func (c *ConfigProvider) GetConfig() *Config {
	return c.config
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnvString("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Type:     getEnvString("DB_TYPE", "inmemory"),
			Host:     getEnvString("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			Name:     getEnvString("DB_NAME", "postgres"),
			User:     getEnvString("DB_USER", "postgres"),
			Password: getEnvString("DB_PASSWORD", "root"),
			SSLMode:  getEnvString("DB_SSL_MODE", "disable"),
		},
		Currency: CurrencyConfig{
			APIKey:  getEnvString("CURRENCY_API_KEY", ""),
			BaseURL: getEnvString("CURRENCY_BASE_URL", "https://api.currencylayer.com"),
			Timeout: getEnvInt("CURRENCY_TIMEOUT", 30),
		},
		Logger: LoggerConfig{
			Level:  getEnvString("LOG_LEVEL", "info"),
			Format: getEnvString("LOG_FORMAT", "json"),
		},
	}
}

// getEnvString gets an environment variable as string with a default value
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an environment variable as int with a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetDatabaseConnectionString returns the database connection string
func (c *ConfigProvider) GetDatabaseConnectionString() string {
	if c.config.Database.Type == "postgres" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.config.Database.Host,
			c.config.Database.Port,
			c.config.Database.User,
			c.config.Database.Password,
			c.config.Database.Name,
			c.config.Database.SSLMode,
		)
	}
	return ""
}

// IsDevelopment returns true if running in development mode
func (c *ConfigProvider) IsDevelopment() bool {
	env := strings.ToLower(getEnvString("ENVIRONMENT", "development"))
	return env == "development" || env == "dev"
}

// IsProduction returns true if running in production mode
func (c *ConfigProvider) IsProduction() bool {
	env := strings.ToLower(getEnvString("ENVIRONMENT", "development"))
	return env == "production" || env == "prod"
}
