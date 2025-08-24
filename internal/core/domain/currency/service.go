package currency

import (
	"fmt"
)

// Currency represents a currency with its information
type Currency struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	IsSupported bool   `json:"is_supported"`
}

// ExchangeRate represents an exchange rate between two currencies
type ExchangeRate struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

// NewCurrency creates a new currency
func NewCurrency(code, name, symbol string, isSupported bool) (*Currency, error) {
	if len(code) != 3 {
		return nil, fmt.Errorf("currency code must be exactly 3 characters")
	}

	if len(name) == 0 {
		return nil, fmt.Errorf("currency name cannot be empty")
	}

	return &Currency{
		Code:        code,
		Name:        name,
		Symbol:      symbol,
		IsSupported: isSupported,
	}, nil
}

// NewExchangeRate creates a new exchange rate
func NewExchangeRate(from, to string, rate float64) (*ExchangeRate, error) {
	if len(from) != 3 || len(to) != 3 {
		return nil, fmt.Errorf("currency codes must be exactly 3 characters")
	}

	if rate <= 0 {
		return nil, fmt.Errorf("exchange rate must be greater than 0")
	}

	return &ExchangeRate{
		From: from,
		To:   to,
		Rate: rate,
	}, nil
}

// Validate validates the currency
func (c *Currency) Validate() error {
	if len(c.Code) != 3 {
		return fmt.Errorf("currency code must be exactly 3 characters")
	}

	if len(c.Name) == 0 {
		return fmt.Errorf("currency name cannot be empty")
	}

	return nil
}

// Validate validates the exchange rate
func (er *ExchangeRate) Validate() error {
	if len(er.From) != 3 || len(er.To) != 3 {
		return fmt.Errorf("currency codes must be exactly 3 characters")
	}

	if er.Rate <= 0 {
		return fmt.Errorf("exchange rate must be greater than 0")
	}

	return nil
}

// CurrencyError represents a currency-related error
type CurrencyError struct {
	Code    string
	Message string
	Cause   error
}

// NewCurrencyError creates a new currency error
func NewCurrencyError(code, message string, cause error) *CurrencyError {
	return &CurrencyError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the error interface
func (e *CurrencyError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CurrencyError) Unwrap() error {
	return e.Cause
}
