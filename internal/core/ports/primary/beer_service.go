package primary

import (
	"context"

	"beers-challenge/internal/core/domain/beers"
)

// BeerService defines the primary port for beer operations
// This represents the use cases from the outside perspective
type BeerService interface {
	CreateBeer(ctx context.Context, req CreateBeerRequest) error
	FindBeerByID(ctx context.Context, id int) (*beers.Beer, error)
	FindAllBeers(ctx context.Context) ([]beers.Beer, error)
	CalculateBoxPrice(ctx context.Context, req CalculateBoxPriceRequest) (*BoxPriceResponse, error)
}

// CreateBeerRequest represents the request to create a beer
type CreateBeerRequest struct {
	ID       int     `json:"id" validate:"required,min=1"`
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Brewery  string  `json:"brewery" validate:"required,min=1,max=100"`
	Country  string  `json:"country" validate:"required,min=1,max=100"`
	Price    float64 `json:"price" validate:"required,min=0"`
	Currency string  `json:"currency" validate:"required,len=3"`
}

// CalculateBoxPriceRequest represents the request to calculate box price
type CalculateBoxPriceRequest struct {
	BeerID   int    `json:"beer_id" validate:"required,min=1"`
	Quantity int    `json:"quantity" validate:"required,min=1,max=1000"`
	Currency string `json:"currency" validate:"required,len=3"`
}

// BoxPriceResponse represents the response for box price calculation
type BoxPriceResponse struct {
	BeerID       int     `json:"beer_id"`
	BeerName     string  `json:"beer_name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	TotalPrice   float64 `json:"total_price"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
}
