package usecases

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"context"
)

// FinderBeersByID is the use case than find all beers
type FinderBeersByID struct {
	beersRepository beers.Repository
}

func NewFinderBeersByID(repository beers.Repository) *FinderBeersByID {
	return &FinderBeersByID{
		repository,
	}
}

// Execute finder a beer by his ID in the repository of beers
func (beersFinder *FinderBeersByID) Execute(ctx context.Context, beerID int) (beers.Beer, error) {
	if err := beers.ValidateBeerID(beerID); err != nil {
		return beers.Beer{}, err
	}
	beerResult, err := beersFinder.beersRepository.FindBeerByID(ctx, beerID)
	if err != nil {
		return beers.Beer{}, err
	}
	return beerResult, nil
}
