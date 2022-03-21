package dependencies

import (
	"os"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/currency"
	"CleverIT-challenge/internal/http/client/currencyLayer"
	"CleverIT-challenge/internal/infrastructure/storage"
)

type Container interface {
	CurrencyService() currency.Service
	BeersRepository() beers.Repository
}

type container struct {
	currencyService currency.Service
	beersRepository beers.Repository
}

const EnvironmentKey = "ENVIRONMENT"

func NewContainer() Container {
	service := currencyLayer.NewCurrencyService()
	environment := os.Getenv(EnvironmentKey)

	return &container{
		currencyService: service,
		beersRepository: storage.New(environment),
	}
}

func (container *container) CurrencyService() currency.Service {
	return container.currencyService
}

func (container *container) BeersRepository() beers.Repository {
	return container.beersRepository
}
