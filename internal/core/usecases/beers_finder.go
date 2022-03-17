package usecases

import (
	"context"

	"CleverIT-challenge/internal/core/domain/beers"
)

// FinderAllBeers is the use case than find all beers
type FinderAllBeers struct {
	beersRepository beers.Repository
}

func NewFinderAllBeers(repository beers.Repository) *FinderAllBeers {
	return &FinderAllBeers{
		repository,
	}
}

// Execute finder in the repository of beers
func (beersFinder *FinderAllBeers) Execute(ctx context.Context) ([]beers.Beer, error) {
	beersList, err := beersFinder.beersRepository.FindAllBeers(ctx)
	if err != nil {
		return nil, err
	}
	return beersList, nil
}
