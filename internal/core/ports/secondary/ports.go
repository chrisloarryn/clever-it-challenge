package secondary

import (
	"context"

	"beers-challenge/internal/core/domain/beers"
)

// BeerRepository defines the secondary port for beer persistence
// This is implemented by the infrastructure layer
type BeerRepository interface {
	Save(ctx context.Context, beer *beers.Beer) error
	FindByID(ctx context.Context, id int) (*beers.Beer, error)
	FindAll(ctx context.Context) ([]beers.Beer, error)
	ExistsByID(ctx context.Context, id int) (bool, error)
}

// CurrencyService defines the secondary port for currency operations
type CurrencyService interface {
	GetExchangeRate(ctx context.Context, from, to string) (float64, error)
	IsValidCurrency(ctx context.Context, currency string) (bool, error)
	GetSupportedCurrencies(ctx context.Context) ([]string, error)
}

// Logger defines the secondary port for logging
type Logger interface {
	Info(ctx context.Context, msg string, fields map[string]interface{})
	Error(ctx context.Context, msg string, err error, fields map[string]interface{})
	Debug(ctx context.Context, msg string, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
}

// ConfigProvider defines the secondary port for configuration
type ConfigProvider interface {
	GetString(key string) string
	GetInt(key string) int
	GetFloat64(key string) float64
	GetBool(key string) bool
}
