package usecases

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"context"
)

// CreateBeer is the use case that create a beer
type CreateBeer struct {
	beersRepository beers.Repository
}

// NewCreateBeer constructor
func NewCreateBeer(repository beers.Repository) *CreateBeer {
	return &CreateBeer{
		repository,
	}
}

// Execute finder in the repository of beers
func (beersCreator *CreateBeer) Execute(ctx context.Context, beer beers.Beer) error {
	if validateError := beers.ValidateBeer(beer); validateError != nil {
		return validateError
	}
	err := beersCreator.beersRepository.SaveBeer(ctx, beer)
	if err != nil {
		return err
	}
	return nil
}
