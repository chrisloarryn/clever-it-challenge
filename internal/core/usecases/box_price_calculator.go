package usecases

import (
	"context"
	"fmt"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/currency"
)

type BoxPriceCalculator struct {
	beersRepository beers.Repository
	currencyService currency.Service
}

func NewBoxPriceCalculator(repository beers.Repository, service currency.Service) *BoxPriceCalculator {
	return &BoxPriceCalculator{
		beersRepository: repository,
		currencyService: service,
	}
}

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
	amountFrom, err = boxPriceCalculator.currencyService.GetCurrencyPriceInDollar(ctx, beerFound.Currency)
	if err != nil {
		return 0, err
	}
	amountTo, err = boxPriceCalculator.currencyService.GetCurrencyPriceInDollar(ctx, currency)
	if err != nil {
		return 0, err
	}
	fromUSD := 1 / amountFrom
	priceTo := fromUSD * amountTo

	return priceTo * float64(quantity), nil
}
