package currencyLayer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"beers-challenge/internal/core/ports/secondary"
)

// CurrencyService implements the secondary.CurrencyService interface
type CurrencyService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// CurrencyLayerResponse represents the API response structure
type CurrencyLayerResponse struct {
	Success   bool               `json:"success"`
	Terms     string             `json:"terms"`
	Privacy   string             `json:"privacy"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
	Error     *CurrencyError     `json:"error,omitempty"`
}

// CurrencyError represents API error response
type CurrencyError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

// NewCurrencyService creates a new CurrencyService instance
func NewCurrencyService(config ConfigProvider) secondary.CurrencyService {
	return &CurrencyService{
		apiKey:  config.GetString("CURRENCY_API_KEY"),
		baseURL: "http://api.currencylayer.com/live",
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ConfigProvider interface for configuration access
type ConfigProvider interface {
	GetString(key string) string
}

// GetExchangeRate gets the exchange rate between two currencies
func (s *CurrencyService) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	// If same currency, rate is 1.0
	if from == to {
		return 1.0, nil
	}

	// CurrencyLayer uses USD as base currency
	// If from is not USD, we need to convert: (1/USD_FROM) * USD_TO
	if from != "USD" && to != "USD" {
		// Get USD/FROM rate
		fromRate, err := s.getUSDRate(ctx, from)
		if err != nil {
			return 0, fmt.Errorf("failed to get rate for %s: %w", from, err)
		}

		// Get USD/TO rate
		toRate, err := s.getUSDRate(ctx, to)
		if err != nil {
			return 0, fmt.Errorf("failed to get rate for %s: %w", to, err)
		}

		// Calculate cross rate: (1/FROM_RATE) * TO_RATE
		return toRate / fromRate, nil
	}

	// Direct conversion with USD
	if from == "USD" {
		return s.getUSDRate(ctx, to)
	}

	// To USD conversion
	fromRate, err := s.getUSDRate(ctx, from)
	if err != nil {
		return 0, err
	}

	return 1.0 / fromRate, nil
}

// getUSDRate gets the USD to target currency rate
func (s *CurrencyService) getUSDRate(ctx context.Context, currency string) (float64, error) {
	if currency == "USD" {
		return 1.0, nil
	}

	url := fmt.Sprintf("%s?access_key=%s&currencies=%s", s.baseURL, s.apiKey, currency)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var response CurrencyLayerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if !response.Success {
		if response.Error != nil {
			return 0, fmt.Errorf("API error: %s (code: %d)", response.Error.Info, response.Error.Code)
		}
		return 0, fmt.Errorf("API request failed")
	}

	// CurrencyLayer returns rates as USD{CURRENCY}
	quoteKey := "USD" + currency
	rate, exists := response.Quotes[quoteKey]
	if !exists {
		return 0, fmt.Errorf("rate not found for currency %s", currency)
	}

	return rate, nil
}

// IsValidCurrency checks if a currency is valid/supported
func (s *CurrencyService) IsValidCurrency(ctx context.Context, currency string) (bool, error) {
	// Common currencies that are typically supported
	supportedCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true,
		"AUD": true, "CAD": true, "CHF": true, "CNY": true,
		"SEK": true, "NZD": true, "MXN": true, "SGD": true,
		"HKD": true, "NOK": true, "TRY": true, "ZAR": true,
		"BRL": true, "INR": true, "RUB": true, "KRW": true,
		"CLP": true, "ARS": true, "COP": true, "PEN": true,
	}

	return supportedCurrencies[currency], nil
}

// GetSupportedCurrencies returns a list of supported currencies
func (s *CurrencyService) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	return []string{
		"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF", "CNY",
		"SEK", "NZD", "MXN", "SGD", "HKD", "NOK", "TRY", "ZAR",
		"BRL", "INR", "RUB", "KRW", "CLP", "ARS", "COP", "PEN",
	}, nil
}
