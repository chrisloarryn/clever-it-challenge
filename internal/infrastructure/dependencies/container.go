package dependencies

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/currency"
)

type Container struct {
	CurrencyService currency.Service
	BeersRepository beers.Repository
}


