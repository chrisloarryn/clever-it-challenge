package usecases

import (
	"context"
	"fmt"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/currency"
)

// BoxPriceCalculator calculates the box price of the beer
type BoxPriceCalculator struct {
	beersRepository beers.Repository
	currencyService currency.Service
}

// NewBoxPriceCalculator constructor
func NewBoxPriceCalculator(repository beers.Repository, service currency.Service) *BoxPriceCalculator {
	return &BoxPriceCalculator{
		beersRepository: repository,
		currencyService: service,
	}
}

// Execute executes the use case BoxPriceCalculator
func (boxPriceCalculator *BoxPriceCalculator) Execute(ctx context.Context, beerID int, quantity int, currency string)(float64, error) {
	if err := beers.ValidateBeerID(beerID); err != nil {
		return 0, err
	}
	if quantity < 1 {
		return 0, fmt.Errorf("invalid quantity")
	}

	beerFound, err := boxPriceCalculator.beersRepository.FindBeerByID(ctx, beerID)
	if err != nil {
		return 0, err
	}
	var amountFrom, amountTo float64
	if amountFrom, err = boxPriceCalculator.currencyService.GetCurrencyPriceInDollar(ctx, beerFound.Currency); err != nil {
		fmt.Println("Error en 1")
		return 0, err
	}
	if amountTo, err = boxPriceCalculator.currencyService.GetCurrencyPriceInDollar(ctx, currency); err != nil {
		fmt.Println("Error en 2")
		return 0, err
	}
	fromUSD := 1 / amountFrom
	priceTo := fromUSD * amountTo

	return priceTo * float64(quantity), nil
}
