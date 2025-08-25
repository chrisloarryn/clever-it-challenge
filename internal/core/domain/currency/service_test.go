package currency

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	usd         = "USD"
	usDollar    = "US Dollar"
	dollarSign  = "$"
	eur         = "EUR"
	testErrCode = "TEST"
	testErr     = "test error"
)

func TestNewCurrency(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c, err := NewCurrency(usd, usDollar, dollarSign, true)
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("invalid code", func(t *testing.T) {
		_, err := NewCurrency("US", usDollar, dollarSign, true)
		assert.Error(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := NewCurrency(usd, "", dollarSign, true)
		assert.Error(t, err)
	})
}

func TestNewExchangeRate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		er, err := NewExchangeRate(usd, eur, 0.85)
		assert.NoError(t, err)
		assert.NotNil(t, er)
	})

	t.Run("invalid from code", func(t *testing.T) {
		_, err := NewExchangeRate("US", eur, 0.85)
		assert.Error(t, err)
	})

	t.Run("invalid to code", func(t *testing.T) {
		_, err := NewExchangeRate(usd, "EU", 0.85)
		assert.Error(t, err)
	})

	t.Run("invalid rate", func(t *testing.T) {
		_, err := NewExchangeRate(usd, eur, 0)
		assert.Error(t, err)
	})
}

func TestCurrencyValidate(t *testing.T) {
	c := &Currency{Code: usd, Name: usDollar}
	assert.NoError(t, c.Validate())

	c.Code = "US"
	assert.Error(t, c.Validate())

	c.Code = usd
	c.Name = ""
	assert.Error(t, c.Validate())
}

func TestExchangeRateValidate(t *testing.T) {
	er := &ExchangeRate{From: usd, To: eur, Rate: 0.85}
	assert.NoError(t, er.Validate())

	er.From = "US"
	assert.Error(t, er.Validate())

	er.From = usd
	er.To = "EU"
	assert.Error(t, er.Validate())

	er.To = eur
	er.Rate = 0
	assert.Error(t, er.Validate())
}

func TestCurrencyError(t *testing.T) {
	err := NewCurrencyError(testErrCode, testErr, nil)
	assert.Equal(t, "TEST: test error", err.Error())
	assert.Nil(t, err.Unwrap())

	cause := errors.New("cause")
	err = NewCurrencyError(testErrCode, testErr, cause)
	assert.Contains(t, err.Error(), "cause")
	assert.Equal(t, cause, err.Unwrap())
}
