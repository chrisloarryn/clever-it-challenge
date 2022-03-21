package beers

import (
	"context"
	"fmt"
)

// Beer represents the beer data
type Beer struct {
	ID       int     `json:"Id"`
	Name     string  `json:"Name"`
	Brewery  string  `json:"Brewery"`
	Country  string  `json:"Country"`
	Price    float64 `json:"Price"`
	Currency string  `json:"Currency"`
}

// ValidateBeerID just validates the ID value shouldn't be negative
func ValidateBeerID(beerID int) error {
	if beerID < 1 {
		return fmt.Errorf("Invalid ID: %d", beerID)
	}
	return nil
}

// ValidatePrice validates the price value
func ValidatePrice(price float64) error {
	if price < 0 {
		return fmt.Errorf("invalid price")
	}
	return nil
}

// ValidateBeer validates all field and required fields of beer data
func ValidateBeer(beer Beer) error {
	if len(beer.Name) == 0 {
		return fmt.Errorf("name couldn't be empty")
	}
	if err := ValidatePrice(beer.Price); err != nil {
		return err
	}
	return nil
}

//go:generate mockgen -package beersmocks -destination beersmocks/beers_repository_mocks.go . Repository

// Repository is the storage abstraction
type Repository interface {
	FindAllBeers(ctx context.Context) ([]Beer, error)
	FindBeerByID(ctx context.Context, beerID int) (Beer, error)
	SaveBeer(ctx context.Context, beer Beer) error
}
