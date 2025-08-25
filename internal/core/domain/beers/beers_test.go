package beers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validID       = 1
	validName     = "Test Beer"
	validBrewery  = "Test Brewery"
	validCountry  = "Chile"
	validPrice    = 1500.0
	validCurrency = "CLP"
	testMessage   = "test message"
)

func TestNewBeerSuccess(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, validName, validBrewery, validCountry, validPrice, validCurrency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, beer)
	assert.Equal(t, validID, beer.ID)
	assert.Equal(t, validName, beer.Name)
	assert.Equal(t, validBrewery, beer.Brewery)
	assert.Equal(t, validCountry, beer.Country)
	assert.Equal(t, validPrice, beer.Price)
	assert.Equal(t, validCurrency, beer.Currency)
	assert.False(t, beer.CreatedAt.IsZero())
	assert.False(t, beer.UpdatedAt.IsZero())
}

func TestNewBeerInvalidID(t *testing.T) {
	// Act
	beer, err := NewBeer(0, validName, validBrewery, validCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "id", validationErr.Field)
}

func TestNewBeerEmptyName(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, "", validBrewery, validCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "name", validationErr.Field)
}

func TestNewBeerLongName(t *testing.T) {
	longName := "This is a very long beer name that exceeds the maximum allowed length of 100 characters and should fail validation"

	// Act
	beer, err := NewBeer(validID, longName, validBrewery, validCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "name", validationErr.Field)
}

func TestNewBeerNegativePrice(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, validName, validBrewery, validCountry, -100.0, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "price", validationErr.Field)
}

func TestNewBeerInvalidCurrency(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, validName, validBrewery, validCountry, validPrice, "INVALID")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "currency", validationErr.Field)
}

func TestNewBeerEmptyBrewery(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, validName, "", validCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "brewery", validationErr.Field)
}

func TestNewBeerLongBrewery(t *testing.T) {
	longBrewery := "This is a very long brewery name that exceeds the maximum allowed length of 100 characters and should fail validation"

	// Act
	beer, err := NewBeer(validID, validName, longBrewery, validCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "brewery", validationErr.Field)
}

func TestNewBeerEmptyCountry(t *testing.T) {
	// Act
	beer, err := NewBeer(validID, validName, validBrewery, "", validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "country", validationErr.Field)
}

func TestNewBeerLongCountry(t *testing.T) {
	longCountry := "This is a very long country name that exceeds the maximum allowed length of 100 characters and should fail validation"

	// Act
	beer, err := NewBeer(validID, validName, validBrewery, longCountry, validPrice, validCurrency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, beer)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "country", validationErr.Field)
}

func TestBeerValidateSuccess(t *testing.T) {
	// Arrange
	beer := &Beer{
		ID:       validID,
		Name:     validName,
		Brewery:  validBrewery,
		Country:  validCountry,
		Price:    validPrice,
		Currency: validCurrency,
	}

	// Act
	err := beer.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestBeerCalculateBoxPriceSuccess(t *testing.T) {
	// Arrange
	beer := &Beer{
		ID:       validID,
		Name:     validName,
		Brewery:  validBrewery,
		Country:  validCountry,
		Price:    validPrice,
		Currency: validCurrency,
	}

	quantity := 24
	exchangeRate := 1.5

	// Act
	totalPrice, err := beer.CalculateBoxPrice(quantity, exchangeRate)

	// Assert
	assert.NoError(t, err)
	expectedPrice := validPrice * float64(quantity) * exchangeRate
	assert.Equal(t, expectedPrice, totalPrice)
}

func TestBeerCalculateBoxPriceInvalidQuantity(t *testing.T) {
	// Arrange
	beer := &Beer{
		ID:       validID,
		Name:     validName,
		Brewery:  validBrewery,
		Country:  validCountry,
		Price:    validPrice,
		Currency: validCurrency,
	}

	// Act
	totalPrice, err := beer.CalculateBoxPrice(0, 1.0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, float64(0), totalPrice)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "quantity", validationErr.Field)
}

func TestBeerCalculateBoxPriceInvalidExchangeRate(t *testing.T) {
	// Arrange
	beer := &Beer{
		ID:       validID,
		Name:     validName,
		Brewery:  validBrewery,
		Country:  validCountry,
		Price:    validPrice,
		Currency: validCurrency,
	}

	// Act
	totalPrice, err := beer.CalculateBoxPrice(24, 0)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, float64(0), totalPrice)
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "exchange_rate", validationErr.Field)
}

func TestBeerGetID(t *testing.T) {
	// Arrange
	beer := &Beer{
		ID:       validID,
		Name:     validName,
		Brewery:  validBrewery,
		Country:  validCountry,
		Price:    validPrice,
		Currency: validCurrency,
	}

	// Act
	beerID := beer.GetID()

	// Assert
	assert.Equal(t, BeerID(validID), beerID)
}

func TestBeerUpdate(t *testing.T) {
	// Arrange
	beer, _ := NewBeer(validID, validName, validBrewery, validCountry, validPrice, validCurrency)
	oldUpdatedAt := beer.UpdatedAt

	// Act
	beer.Update()

	// Assert
	assert.NotEqual(t, oldUpdatedAt, beer.UpdatedAt)
}

func TestValidationErrorError(t *testing.T) {
	// Arrange
	validationErr := NewValidationError("test_field", testMessage)

	// Act
	errorMsg := validationErr.Error()

	// Assert
	assert.Contains(t, errorMsg, "test_field")
	assert.Contains(t, errorMsg, testMessage)
	assert.Contains(t, errorMsg, "validation error")
}

func TestDomainErrorError(t *testing.T) {
	// Arrange
	domainErr := NewDomainError("TEST_ERROR", testMessage, nil)

	// Act
	errorMsg := domainErr.Error()

	// Assert
	assert.Contains(t, errorMsg, "TEST_ERROR")
	assert.Contains(t, errorMsg, testMessage)
}

func TestDomainErrorWithCause(t *testing.T) {
	// Arrange
	cause := assert.AnError
	domainErr := NewDomainError("TEST_ERROR", testMessage, cause)

	// Act
	errorMsg := domainErr.Error()
	unwrappedErr := domainErr.Unwrap()

	// Assert
	assert.Contains(t, errorMsg, "TEST_ERROR")
	assert.Contains(t, errorMsg, testMessage)
	assert.Contains(t, errorMsg, "caused by")
	assert.Equal(t, cause, unwrappedErr)
}
