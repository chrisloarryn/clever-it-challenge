package beers

import "context"

type Beer struct {
	ID int `json:"Id"`
	Name string `json:"Name"`
	Brewery string `json:"Brewery"`
	Country string `json:"Country"`
	Price float64 `json:"Price"`
	Currency string `json:"Currency"`
}

//go:generate mockgen -package beersmocks -destination beersmocks/beers_repository_mocks.go . Repository

type Repository interface {
	FindAllBeers(ctx context.Context) ([]Beer, error)
	FindBeerByID(ctx context.Context, beerID int) (Beer, error)
}