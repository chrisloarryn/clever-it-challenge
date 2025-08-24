package currencyLayer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"beers-challenge/internal/core/domain/currency"
	"beers-challenge/internal/core/ports/secondary"
	"beers-challenge/internal/infrastructure/config"
)

// Client implements the secondary.CurrencyService interface
type Client struct {
	config     *config.ConfigProvider
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewCurrencyService creates a new currency service client
func NewCurrencyService(config *config.ConfigProvider) secondary.CurrencyService {
	timeout := time.Duration(config.GetInt("currency.timeout")) * time.Second

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: config.GetString("currency.base_url"),
		apiKey:  config.GetString("currency.api_key"),
	}
}

// Response structures for CurrencyLayer API
type CurrencyLayerResponse struct {
	Success   bool                `json:"success"`
	Error     *CurrencyLayerError `json:"error,omitempty"`
	Quotes    map[string]float64  `json:"quotes,omitempty"`
	Source    string              `json:"source,omitempty"`
	Timestamp int64               `json:"timestamp,omitempty"`
}

type CurrencyLayerError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

type SupportedCurrenciesResponse struct {
	Success    bool                `json:"success"`
	Error      *CurrencyLayerError `json:"error,omitempty"`
	Currencies map[string]string   `json:"currencies,omitempty"`
}

// GetExchangeRate gets the exchange rate between two currencies
func (c *Client) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}

	// CurrencyLayer uses USD as base currency
	if from != "USD" && to != "USD" {
		// Convert from -> USD -> to
		fromToUSD, err := c.getRate(ctx, from, "USD")
		if err != nil {
			return 0, err
		}

		usdToTo, err := c.getRate(ctx, "USD", to)
		if err != nil {
			return 0, err
		}

		return fromToUSD * usdToTo, nil
	}

	return c.getRate(ctx, from, to)
}

// getRate gets a single exchange rate
func (c *Client) getRate(ctx context.Context, from, to string) (float64, error) {
	if c.apiKey == "" {
		// Return mock rate for development/testing
		return c.getMockRate(from, to), nil
	}

	url := fmt.Sprintf("%s/live?access_key=%s&currencies=%s&source=%s",
		c.baseURL, c.apiKey, to, from)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, currency.NewCurrencyError("REQUEST_CREATION_FAILED", "Failed to create request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, currency.NewCurrencyError("API_REQUEST_FAILED", "Failed to make API request", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, currency.NewCurrencyError("RESPONSE_READ_FAILED", "Failed to read response", err)
	}

	var response CurrencyLayerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, currency.NewCurrencyError("JSON_DECODE_FAILED", "Failed to decode response", err)
	}

	if !response.Success {
		if response.Error != nil {
			return 0, currency.NewCurrencyError("API_ERROR", response.Error.Info, nil)
		}
		return 0, currency.NewCurrencyError("API_ERROR", "Unknown API error", nil)
	}

	quoteKey := from + to
	rate, exists := response.Quotes[quoteKey]
	if !exists {
		return 0, currency.NewCurrencyError("RATE_NOT_FOUND",
			fmt.Sprintf("Exchange rate not found for %s to %s", from, to), nil)
	}

	return rate, nil
}

// IsValidCurrency checks if a currency is valid
func (c *Client) IsValidCurrency(ctx context.Context, currencyCode string) (bool, error) {
	currencyCode = strings.ToUpper(strings.TrimSpace(currencyCode))

	if len(currencyCode) != 3 {
		return false, nil
	}

	// Common currencies for quick validation
	commonCurrencies := map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true,
		"AUD": true, "CAD": true, "CHF": true, "CNY": true,
		"SEK": true, "NZD": true, "MXN": true, "SGD": true,
		"HKD": true, "NOK": true, "TRY": true, "ZAR": true,
		"BRL": true, "INR": true, "KRW": true, "RUB": true,
		"CLP": true, "PEN": true, "COP": true, "ARS": true,
	}

	if commonCurrencies[currencyCode] {
		return true, nil
	}

	// If API key is available, check with API
	if c.apiKey != "" {
		return c.validateCurrencyWithAPI(ctx, currencyCode)
	}

	// For development without API key, accept any 3-letter code
	return true, nil
}

// validateCurrencyWithAPI validates currency using the API
func (c *Client) validateCurrencyWithAPI(ctx context.Context, currencyCode string) (bool, error) {
	supportedCurrencies, err := c.GetSupportedCurrencies(ctx)
	if err != nil {
		// If we can't get supported currencies, fall back to common currencies
		return false, nil
	}

	for _, supportedCode := range supportedCurrencies {
		if supportedCode == currencyCode {
			return true, nil
		}
	}

	return false, nil
}

// GetSupportedCurrencies gets all supported currencies
func (c *Client) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	if c.apiKey == "" {
		// Return common currencies for development/testing
		return []string{
			"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF", "CNY",
			"SEK", "NZD", "MXN", "SGD", "HKD", "NOK", "TRY", "ZAR",
			"BRL", "INR", "KRW", "RUB", "CLP", "PEN", "COP", "ARS",
		}, nil
	}

	url := fmt.Sprintf("%s/list?access_key=%s", c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, currency.NewCurrencyError("REQUEST_CREATION_FAILED", "Failed to create request", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, currency.NewCurrencyError("API_REQUEST_FAILED", "Failed to make API request", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, currency.NewCurrencyError("RESPONSE_READ_FAILED", "Failed to read response", err)
	}

	var response SupportedCurrenciesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, currency.NewCurrencyError("JSON_DECODE_FAILED", "Failed to decode response", err)
	}

	if !response.Success {
		if response.Error != nil {
			return nil, currency.NewCurrencyError("API_ERROR", response.Error.Info, nil)
		}
		return nil, currency.NewCurrencyError("API_ERROR", "Unknown API error", nil)
	}

	var currencies []string
	for code := range response.Currencies {
		currencies = append(currencies, code)
	}

	return currencies, nil
}

// getMockRate returns a mock exchange rate for development/testing
func (c *Client) getMockRate(from, to string) float64 {
	// Mock exchange rates for testing
	rates := map[string]map[string]float64{
		"USD": {
			"EUR": 0.85, "GBP": 0.73, "JPY": 110.0, "CAD": 1.25,
			"AUD": 1.35, "CHF": 0.92, "CNY": 6.45, "CLP": 800.0,
		},
		"EUR": {
			"USD": 1.18, "GBP": 0.86, "JPY": 129.0, "CLP": 940.0,
		},
		"CLP": {
			"USD": 0.00125, "EUR": 0.00106, "GBP": 0.00091,
		},
	}

	if fromRates, exists := rates[from]; exists {
		if rate, exists := fromRates[to]; exists {
			return rate
		}
	}

	// If no specific rate found, return a default
	return 1.0
}
