package beers

import (
	"fmt"
	"strings"
	"time"
)

const (
	// Validation error messages
	ErrMustBeGreaterThanZero = "must be greater than 0"
	ErrCannotBeEmpty         = "cannot be empty"
	ErrCannotExceed100Chars  = "cannot exceed 100 characters"
	ErrCannotBeNegative      = "cannot be negative"
	ErrMustBe3Characters     = "must be exactly 3 characters (ISO 4217)"
)

// Beer represents the beer domain entity
type Beer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Brewery   string    `json:"brewery"`
	Country   string    `json:"country"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeerID represents a beer identifier
type BeerID int

// NewBeer creates a new beer with validation
func NewBeer(id int, name, brewery, country string, price float64, currency string) (*Beer, error) {
	beer := &Beer{
		ID:        id,
		Name:      strings.TrimSpace(name),
		Brewery:   strings.TrimSpace(brewery),
		Country:   strings.TrimSpace(country),
		Price:     price,
		Currency:  strings.ToUpper(strings.TrimSpace(currency)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := beer.Validate(); err != nil {
		return nil, err
	}

	return beer, nil
}

// Validate validates the beer entity
func (b *Beer) Validate() error {
	if b.ID < 1 {
		return NewValidationError("id", ErrMustBeGreaterThanZero)
	}

	if len(b.Name) == 0 {
		return NewValidationError("name", ErrCannotBeEmpty)
	}

	if len(b.Name) > 100 {
		return NewValidationError("name", ErrCannotExceed100Chars)
	}

	if len(b.Brewery) == 0 {
		return NewValidationError("brewery", ErrCannotBeEmpty)
	}

	if len(b.Brewery) > 100 {
		return NewValidationError("brewery", ErrCannotExceed100Chars)
	}

	if len(b.Country) == 0 {
		return NewValidationError("country", ErrCannotBeEmpty)
	}

	if len(b.Country) > 100 {
		return NewValidationError("country", ErrCannotExceed100Chars)
	}

	if b.Price < 0 {
		return NewValidationError("price", ErrCannotBeNegative)
	}

	if len(b.Currency) != 3 {
		return NewValidationError("currency", ErrMustBe3Characters)
	}

	return nil
}

// CalculateBoxPrice calculates the total price for a quantity of beers
func (b *Beer) CalculateBoxPrice(quantity int, exchangeRate float64) (float64, error) {
	if quantity < 1 {
		return 0, NewValidationError("quantity", ErrMustBeGreaterThanZero)
	}

	if exchangeRate <= 0 {
		return 0, NewValidationError("exchange_rate", ErrMustBeGreaterThanZero)
	}

	return b.Price * float64(quantity) * exchangeRate, nil
}

// Update updates the beer's updatedAt timestamp
func (b *Beer) Update() {
	b.UpdatedAt = time.Now()
}

// GetID returns the beer ID
func (b *Beer) GetID() BeerID {
	return BeerID(b.ID)
}

// ValidationError represents a domain validation error
type ValidationError struct {
	Field   string
	Message string
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Cause   error
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *DomainError) Unwrap() error {
	return e.Cause
}
