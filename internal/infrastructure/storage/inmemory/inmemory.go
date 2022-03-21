package inmemory

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"context"
	"fmt"
)

type InMemoryRepository struct {
	list map[int]beers.Beer
}

func (repository *InMemoryRepository) FindAllBeers(_ context.Context) ([]beers.Beer, error) {
	result := []beers.Beer{}

	for _, beer := range repository.list {
		result = append(result, beer)
	}

	return result, nil
}

func (repository *InMemoryRepository) FindBeerByID(_ context.Context, beerID int) (beers.Beer, error) {

	for key, beer := range repository.list {
		if key == beerID {
			return beer, nil
		}
	}
	return beers.Beer{}, fmt.Errorf("Beer ID doesn't exist")
}

func (repository *InMemoryRepository) SaveBeer(_ context.Context, beer beers.Beer) error {
	_, exist := repository.list[beer.ID]
	if exist {
		return fmt.Errorf("The beer ID already exists")
	}
	repository.list[beer.ID] = beer
	return nil
}

func NewInMemoryRepository() beers.Repository {
	return &InMemoryRepository{
		list: map[int]beers.Beer{},
	}
}
