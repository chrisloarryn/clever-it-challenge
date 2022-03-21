package inmemory

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"context"
	"fmt"
)

// Repository is the struct when you choose the in memory storage
type Repository struct {
	list map[int]beers.Beer
}

func (repository *Repository) FindAllBeers(_ context.Context) ([]beers.Beer, error) {
	var result []beers.Beer

	for _, beer := range repository.list {
		result = append(result, beer)
	}

	return result, nil
}

func (repository *Repository) FindBeerByID(_ context.Context, beerID int) (beers.Beer, error) {

	for key, beer := range repository.list {
		if key == beerID {
			return beer, nil
		}
	}
	return beers.Beer{}, fmt.Errorf("beer ID doesn't exist")
}

func (repository *Repository) SaveBeer(_ context.Context, beer beers.Beer) error {
	_, exist := repository.list[beer.ID]
	if exist {
		return fmt.Errorf("the beer ID already exists")
	}
	repository.list[beer.ID] = beer
	return nil
}

func NewInMemoryRepository() beers.Repository {
	return &Repository{
		list: map[int]beers.Beer{},
	}
}
