package usecases

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/currency"
	"context"
	"fmt"
)

// CreateBeer is the use case that create a beer
type CreateBeer struct {
	beersRepository beers.Repository
	currencyService currency.Service
}

// NewCreateBeer constructor
func NewCreateBeer(repository beers.Repository, service currency.Service) *CreateBeer {
	return &CreateBeer{
		repository,
		service,
	}
}

// Execute finder in the repository of beers
func (beersCreator *CreateBeer) Execute(ctx context.Context, beer beers.Beer) error {
	if validateError := beers.ValidateBeer(beer); validateError != nil {
		return validateError
	}
	if isValid, err := beersCreator.currencyService.IsValidCurrency(ctx, beer.Currency); err != nil {
		return err
	} else if !isValid {
		return fmt.Errorf("invalid currency value")
	}
	err := beersCreator.beersRepository.SaveBeer(ctx, beer)
	if err != nil {
		return err
	}
	return nil
}
