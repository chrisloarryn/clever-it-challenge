package currencyLayer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"beers-challenge/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
)

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
	currencyBaseURL   = "currency.base_url"
	currencyAPIKey    = "currency.api_key"
)

func setupTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestGetExchangeRate(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		fmt.Fprintln(w, `{"success": true, "quotes": {"USDCLP": 800.0}}`)
	})
	defer server.Close()

	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.BaseURL = server.URL
	cfg.GetConfig().Currency.APIKey = "test_key"

	client := NewCurrencyService(cfg)

	rate, err := client.GetExchangeRate(context.Background(), "USD", "CLP")
	assert.NoError(t, err)
	assert.Equal(t, 800.0, rate)
}

func TestIsValidCurrency(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		fmt.Fprintln(w, `{"success": true, "currencies": {"CLP": "Chilean Peso"}}`)
	})
	defer server.Close()

	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.BaseURL = server.URL
	cfg.GetConfig().Currency.APIKey = "test_key"

	client := NewCurrencyService(cfg)

	valid, err := client.IsValidCurrency(context.Background(), "CLP")
	assert.NoError(t, err)
	assert.True(t, valid)

	valid, err = client.IsValidCurrency(context.Background(), "INVALID")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestGetSupportedCurrencies(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		fmt.Fprintln(w, `{"success": true, "currencies": {"USD": "United States Dollar", "CLP": "Chilean Peso"}}`)
	})
	defer server.Close()

	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.BaseURL = server.URL
	cfg.GetConfig().Currency.APIKey = "test_key"

	client := NewCurrencyService(cfg)

	currencies, err := client.GetSupportedCurrencies(context.Background())
	assert.NoError(t, err)
	assert.Contains(t, currencies, "USD")
	assert.Contains(t, currencies, "CLP")
}

func TestGetExchangeRateNoAPIKey(t *testing.T) {
	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.APIKey = ""

	client := NewCurrencyService(cfg)

	rate, err := client.GetExchangeRate(context.Background(), "USD", "EUR")
	assert.NoError(t, err)
	assert.Equal(t, 0.85, rate) // Mocked rate
}

func TestGetExchangeRateComplexConversion(t *testing.T) {
	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.APIKey = ""

	client := NewCurrencyService(cfg)

	rate, err := client.GetExchangeRate(context.Background(), "EUR", "CLP")
	assert.NoError(t, err)
	assert.InDelta(t, 940.0, rate, 0.001) // Mocked rate
}

func TestIsValidCurrencyShortCode(t *testing.T) {
	client := NewCurrencyService(config.NewConfigProvider())
	valid, err := client.IsValidCurrency(context.Background(), "US")
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestGetRateAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		fmt.Fprintln(w, `{"success": false, "error": {"code": 101, "info": "api error"}}`)
	}))
	defer server.Close()

	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.BaseURL = server.URL
	cfg.GetConfig().Currency.APIKey = "test_key"

	client := NewCurrencyService(cfg)
	_, err := client.GetExchangeRate(context.Background(), "USD", "CLP")
	assert.Error(t, err)
}

func TestGetRateRateNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		fmt.Fprintln(w, `{"success": true, "quotes": {}}`)
	}))
	defer server.Close()

	cfg := config.NewConfigProvider()
	cfg.GetConfig().Currency.BaseURL = server.URL
	cfg.GetConfig().Currency.APIKey = "test_key"

	client := NewCurrencyService(cfg)
	_, err := client.GetExchangeRate(context.Background(), "USD", "CLP")
	assert.Error(t, err)
}
