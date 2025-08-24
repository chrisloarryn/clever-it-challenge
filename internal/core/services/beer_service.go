package services

import (
	"context"
	"fmt"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/core/ports/secondary"
)

// BeerServiceImpl implements the BeerService primary port
type BeerServiceImpl struct {
	beerRepo        secondary.BeerRepository
	currencyService secondary.CurrencyService
	logger          secondary.Logger
}

// NewBeerService creates a new beer service
func NewBeerService(
	beerRepo secondary.BeerRepository,
	currencyService secondary.CurrencyService,
	logger secondary.Logger,
) primary.BeerService {
	return &BeerServiceImpl{
		beerRepo:        beerRepo,
		currencyService: currencyService,
		logger:          logger,
	}
}

// CreateBeer creates a new beer
func (s *BeerServiceImpl) CreateBeer(ctx context.Context, req primary.CreateBeerRequest) error {
	s.logger.Info(ctx, "Creating beer", map[string]interface{}{
		"beer_id": req.ID,
		"name":    req.Name,
	})

	// Check if beer already exists
	exists, err := s.beerRepo.ExistsByID(ctx, req.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to check beer existence", err, map[string]interface{}{
			"beer_id": req.ID,
		})
		return fmt.Errorf("failed to check beer existence: %w", err)
	}

	if exists {
		return beers.NewDomainError("BEER_ALREADY_EXISTS", "Beer with this ID already exists", nil)
	}

	// Validate currency
	isValid, err := s.currencyService.IsValidCurrency(ctx, req.Currency)
	if err != nil {
		s.logger.Error(ctx, "Failed to validate currency", err, map[string]interface{}{
			"currency": req.Currency,
		})
		return fmt.Errorf("failed to validate currency: %w", err)
	}

	if !isValid {
		return beers.NewDomainError("INVALID_CURRENCY", "Invalid currency code", nil)
	}

	// Create domain entity
	beer, err := beers.NewBeer(req.ID, req.Name, req.Brewery, req.Country, req.Price, req.Currency)
	if err != nil {
		s.logger.Error(ctx, "Failed to create beer entity", err, map[string]interface{}{
			"beer_id": req.ID,
		})
		return fmt.Errorf("failed to create beer: %w", err)
	}

	// Save beer
	if err := s.beerRepo.Save(ctx, beer); err != nil {
		s.logger.Error(ctx, "Failed to save beer", err, map[string]interface{}{
			"beer_id": req.ID,
		})
		return fmt.Errorf("failed to save beer: %w", err)
	}

	s.logger.Info(ctx, "Beer created successfully", map[string]interface{}{
		"beer_id": req.ID,
	})

	return nil
}

// FindBeerByID finds a beer by its ID
func (s *BeerServiceImpl) FindBeerByID(ctx context.Context, id int) (*beers.Beer, error) {
	s.logger.Debug(ctx, "Finding beer by ID", map[string]interface{}{
		"beer_id": id,
	})

	if id < 1 {
		return nil, beers.NewValidationError("id", beers.ErrMustBeGreaterThanZero)
	}

	beer, err := s.beerRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.Error(ctx, "Failed to find beer", err, map[string]interface{}{
			"beer_id": id,
		})
		return nil, fmt.Errorf("failed to find beer: %w", err)
	}

	return beer, nil
}

// FindAllBeers finds all beers
func (s *BeerServiceImpl) FindAllBeers(ctx context.Context) ([]beers.Beer, error) {
	s.logger.Debug(ctx, "Finding all beers", nil)

	beersSlice, err := s.beerRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to find all beers", err, nil)
		return nil, fmt.Errorf("failed to find all beers: %w", err)
	}

	s.logger.Info(ctx, "Found beers", map[string]interface{}{
		"count": len(beersSlice),
	})

	return beersSlice, nil
}

// CalculateBoxPrice calculates the price for a box of beers
func (s *BeerServiceImpl) CalculateBoxPrice(ctx context.Context, req primary.CalculateBoxPriceRequest) (*primary.BoxPriceResponse, error) {
	s.logger.Info(ctx, "Calculating box price", map[string]interface{}{
		"beer_id":  req.BeerID,
		"quantity": req.Quantity,
		"currency": req.Currency,
	})

	// Find the beer
	beer, err := s.FindBeerByID(ctx, req.BeerID)
	if err != nil {
		return nil, err
	}

	// Get exchange rate
	exchangeRate := 1.0
	if beer.Currency != req.Currency {
		rate, err := s.currencyService.GetExchangeRate(ctx, beer.Currency, req.Currency)
		if err != nil {
			s.logger.Error(ctx, "Failed to get exchange rate", err, map[string]interface{}{
				"from": beer.Currency,
				"to":   req.Currency,
			})
			return nil, fmt.Errorf("failed to get exchange rate: %w", err)
		}
		exchangeRate = rate
	}

	// Calculate total price
	totalPrice, err := beer.CalculateBoxPrice(req.Quantity, exchangeRate)
	if err != nil {
		s.logger.Error(ctx, "Failed to calculate box price", err, map[string]interface{}{
			"beer_id":       req.BeerID,
			"quantity":      req.Quantity,
			"exchange_rate": exchangeRate,
		})
		return nil, fmt.Errorf("failed to calculate box price: %w", err)
	}

	response := &primary.BoxPriceResponse{
		BeerID:     req.BeerID,
		BeerName:   beer.Name,
		Quantity:   req.Quantity,
		UnitPrice:  beer.Price * exchangeRate,
		TotalPrice: totalPrice,
		Currency:   req.Currency,
	}

	if beer.Currency != req.Currency {
		response.ExchangeRate = exchangeRate
	}

	s.logger.Info(ctx, "Box price calculated successfully", map[string]interface{}{
		"beer_id":     req.BeerID,
		"total_price": totalPrice,
		"currency":    req.Currency,
	})

	return response, nil
}
