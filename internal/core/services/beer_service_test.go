package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/infrastructure/logger"
)

const (
	testBeerName = "Test Beer"
	testBrewery  = "Test Brewery"
	testCountry  = "Chile"
	testCurrency = "CLP"
	testPrice    = 1500.0
	testBeerID   = 1
)

// Mock implementations
type MockBeerRepository struct {
	mock.Mock
}

func (m *MockBeerRepository) Save(ctx context.Context, beer *beers.Beer) error {
	args := m.Called(ctx, beer)
	return args.Error(0)
}

func (m *MockBeerRepository) FindByID(ctx context.Context, id int) (*beers.Beer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*beers.Beer), args.Error(1)
}

func (m *MockBeerRepository) FindAll(ctx context.Context) ([]beers.Beer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]beers.Beer), args.Error(1)
}

func (m *MockBeerRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockCurrencyService struct {
	mock.Mock
}

func (m *MockCurrencyService) GetExchangeRate(ctx context.Context, from, to string) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockCurrencyService) IsValidCurrency(ctx context.Context, currency string) (bool, error) {
	args := m.Called(ctx, currency)
	return args.Bool(0), args.Error(1)
}

func (m *MockCurrencyService) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateBeerSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	req := primary.CreateBeerRequest{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("ExistsByID", ctx, testBeerID).Return(false, nil)
	mockCurrency.On("IsValidCurrency", ctx, testCurrency).Return(true, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*beers.Beer")).Return(nil)

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockCurrency.AssertExpectations(t)
}

func TestCreateBeerAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	req := primary.CreateBeerRequest{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("ExistsByID", ctx, testBeerID).Return(true, nil)

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
	domainErr, ok := err.(*beers.DomainError)
	assert.True(t, ok)
	assert.Equal(t, "BEER_ALREADY_EXISTS", domainErr.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateBeerInvalidCurrency(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	req := primary.CreateBeerRequest{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: "XXX",
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("ExistsByID", ctx, testBeerID).Return(false, nil)
	mockCurrency.On("IsValidCurrency", ctx, "XXX").Return(false, nil)

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
	domainErr, ok := err.(*beers.DomainError)
	assert.True(t, ok)
	assert.Equal(t, "INVALID_CURRENCY", domainErr.Code)
	mockRepo.AssertExpectations(t)
	mockCurrency.AssertExpectations(t)
}

func TestFindBeerByIDSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	expectedBeer := &beers.Beer{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("FindByID", ctx, testBeerID).Return(expectedBeer, nil)

	// Act
	result, err := service.FindBeerByID(ctx, testBeerID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedBeer, result)
	mockRepo.AssertExpectations(t)
}

func TestFindBeerByIDInvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	ctx := context.Background()

	// Act
	result, err := service.FindBeerByID(ctx, 0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	validationErr, ok := err.(*beers.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "id", validationErr.Field)
}

func TestFindBeerByIDNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	ctx := context.Background()
	notFoundErr := beers.NewDomainError("BEER_NOT_FOUND", "Beer not found", nil)

	// Setup mocks
	mockRepo.On("FindByID", ctx, 999).Return(nil, notFoundErr)

	// Act
	result, err := service.FindBeerByID(ctx, 999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find beer")
	mockRepo.AssertExpectations(t)
}

func TestFindAllBeersSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	expectedBeers := []beers.Beer{
		{ID: 1, Name: "Beer 1"},
		{ID: 2, Name: "Beer 2"},
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("FindAll", ctx).Return(expectedBeers, nil)

	// Act
	result, err := service.FindAllBeers(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedBeers, result)
	mockRepo.AssertExpectations(t)
}

func TestFindAllBeersError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	ctx := context.Background()
	expectedErr := errors.New("database error")

	// Setup mocks
	mockRepo.On("FindAll", ctx).Return([]beers.Beer{}, expectedErr)

	// Act
	result, err := service.FindAllBeers(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), expectedErr.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateBeerExistsError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()
	service := NewBeerService(mockRepo, mockCurrency, logger)
	req := primary.CreateBeerRequest{ID: 1}
	ctx := context.Background()
	mockRepo.On("ExistsByID", ctx, 1).Return(false, errors.New("db error"))

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
}

func TestCreateBeerInvalidCurrencyError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()
	service := NewBeerService(mockRepo, mockCurrency, logger)
	req := primary.CreateBeerRequest{ID: 1, Currency: "XXX"}
	ctx := context.Background()
	mockRepo.On("ExistsByID", ctx, 1).Return(false, nil)
	mockCurrency.On("IsValidCurrency", ctx, "XXX").Return(false, errors.New("currency service error"))

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
}

func TestCreateBeerNewBeerError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()
	service := NewBeerService(mockRepo, mockCurrency, logger)
	req := primary.CreateBeerRequest{ID: 0} // Invalid ID
	ctx := context.Background()
	mockRepo.On("ExistsByID", ctx, 0).Return(false, nil)
	mockCurrency.On("IsValidCurrency", ctx, "").Return(true, nil)

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
}

func TestCreateBeerSaveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()
	service := NewBeerService(mockRepo, mockCurrency, logger)
	req := primary.CreateBeerRequest{ID: 1, Name: "Test", Brewery: "Test", Country: "Test", Price: 1, Currency: "USD"}
	ctx := context.Background()
	mockRepo.On("ExistsByID", ctx, 1).Return(false, nil)
	mockCurrency.On("IsValidCurrency", ctx, "USD").Return(true, nil)
	mockRepo.On("Save", ctx, mock.Anything).Return(errors.New("db error"))

	// Act
	err := service.CreateBeer(ctx, req)

	// Assert
	assert.Error(t, err)
}

func TestCalculateBoxPriceCalculateError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()
	service := NewBeerService(mockRepo, mockCurrency, logger)
	beer := &beers.Beer{ID: 1, Name: "Test", Brewery: "Test", Country: "Test", Price: 1, Currency: "USD"}
	req := primary.CalculateBoxPriceRequest{BeerID: 1, Quantity: 0} // Invalid quantity
	ctx := context.Background()
	mockRepo.On("FindByID", ctx, 1).Return(beer, nil)

	// Act
	_, err := service.CalculateBoxPrice(ctx, req)

	// Assert
	assert.Error(t, err)
}

func TestCalculateBoxPriceSuccess(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	beer := &beers.Beer{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	req := primary.CalculateBoxPriceRequest{
		BeerID:   testBeerID,
		Quantity: 24,
		Currency: "USD",
	}

	ctx := context.Background()
	exchangeRate := 0.00125 // CLP to USD

	// Setup mocks
	mockRepo.On("FindByID", ctx, testBeerID).Return(beer, nil)
	mockCurrency.On("GetExchangeRate", ctx, testCurrency, "USD").Return(exchangeRate, nil)

	// Act
	result, err := service.CalculateBoxPrice(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testBeerID, result.BeerID)
	assert.Equal(t, testBeerName, result.BeerName)
	assert.Equal(t, 24, result.Quantity)
	assert.Equal(t, "USD", result.Currency)
	assert.Equal(t, exchangeRate, result.ExchangeRate)

	expectedUnitPrice := beer.Price * exchangeRate
	expectedTotalPrice := expectedUnitPrice * float64(req.Quantity)
	assert.Equal(t, expectedUnitPrice, result.UnitPrice)
	assert.Equal(t, expectedTotalPrice, result.TotalPrice)

	mockRepo.AssertExpectations(t)
	mockCurrency.AssertExpectations(t)
}

func TestCalculateBoxPriceSameCurrency(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	beer := &beers.Beer{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	req := primary.CalculateBoxPriceRequest{
		BeerID:   testBeerID,
		Quantity: 12,
		Currency: testCurrency, // Same currency as beer
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("FindByID", ctx, testBeerID).Return(beer, nil)
	// No currency service call expected since same currency

	// Act
	result, err := service.CalculateBoxPrice(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testBeerID, result.BeerID)
	assert.Equal(t, testBeerName, result.BeerName)
	assert.Equal(t, 12, result.Quantity)
	assert.Equal(t, testCurrency, result.Currency)
	assert.Equal(t, float64(0), result.ExchangeRate) // No exchange rate for same currency

	expectedUnitPrice := beer.Price
	expectedTotalPrice := expectedUnitPrice * float64(req.Quantity)
	assert.Equal(t, expectedUnitPrice, result.UnitPrice)
	assert.Equal(t, expectedTotalPrice, result.TotalPrice)

	mockRepo.AssertExpectations(t)
	mockCurrency.AssertExpectations(t)
}

func TestCalculateBoxPriceExchangeRateError(t *testing.T) {
	// Arrange
	mockRepo := new(MockBeerRepository)
	mockCurrency := new(MockCurrencyService)
	logger := logger.NewNoOpLogger()

	service := NewBeerService(mockRepo, mockCurrency, logger)

	beer := &beers.Beer{
		ID:       testBeerID,
		Name:     testBeerName,
		Brewery:  testBrewery,
		Country:  testCountry,
		Price:    testPrice,
		Currency: testCurrency,
	}

	req := primary.CalculateBoxPriceRequest{
		BeerID:   testBeerID,
		Quantity: 24,
		Currency: "USD",
	}

	ctx := context.Background()

	// Setup mocks
	mockRepo.On("FindByID", ctx, testBeerID).Return(beer, nil)
	mockCurrency.On("GetExchangeRate", ctx, testCurrency, "USD").Return(0.0, errors.New("currency service error"))

	// Act
	result, err := service.CalculateBoxPrice(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get exchange rate")
	mockRepo.AssertExpectations(t)
	mockCurrency.AssertExpectations(t)
}
