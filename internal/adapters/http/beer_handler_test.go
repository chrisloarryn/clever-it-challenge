package http

import (
	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/infrastructure/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	beersEndpoint     = "/beers"
	testBeerName      = "Test Beer"
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
	serviceErr        = "service error"
)

// MockBeerService is a mock of BeerService
type MockBeerService struct {
	mock.Mock
}

func (m *MockBeerService) CreateBeer(ctx context.Context, req primary.CreateBeerRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockBeerService) FindBeerByID(ctx context.Context, id int) (*beers.Beer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*beers.Beer), args.Error(1)
}

func (m *MockBeerService) FindAllBeers(ctx context.Context) ([]beers.Beer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]beers.Beer), args.Error(1)
}

func (m *MockBeerService) CalculateBoxPrice(ctx context.Context, req primary.CalculateBoxPriceRequest) (*primary.BoxPriceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*primary.BoxPriceResponse), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestCreateBeer(t *testing.T) {
	mockService := new(MockBeerService)
	log := logger.NewNoOpLogger()
	handler := NewBeerHandler(mockService, log)

	r := setupRouter()
	r.POST(beersEndpoint, handler.CreateBeer)

	t.Run("success", func(t *testing.T) {
		reqBody := primary.CreateBeerRequest{
			ID: 1, Name: testBeerName, Brewery: "Test", Country: "Test", Price: 1.0, Currency: "USD",
		}
		mockService.On("CreateBeer", mock.Anything, reqBody).Return(nil).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, beersEndpoint, bytes.NewBuffer(body))
		req.Header.Set(contentTypeHeader, jsonContentType)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, beersEndpoint, bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set(contentTypeHeader, jsonContentType)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run(serviceErr, func(t *testing.T) {
		reqBody := primary.CreateBeerRequest{
			ID: 1, Name: testBeerName, Brewery: "Test", Country: "Test", Price: 1.0, Currency: "USD",
		}
		mockService.On("CreateBeer", mock.Anything, reqBody).Return(errors.New(serviceErr)).Once()

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, beersEndpoint, bytes.NewBuffer(body))
		req.Header.Set(contentTypeHeader, jsonContentType)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetBeer(t *testing.T) {
	mockService := new(MockBeerService)
	log := logger.NewNoOpLogger()
	handler := NewBeerHandler(mockService, log)

	r := setupRouter()
	r.GET("/beers/:id", handler.GetBeer)

	t.Run("success", func(t *testing.T) {
		beer := &beers.Beer{ID: 1, Name: testBeerName}
		mockService.On("FindBeerByID", mock.Anything, 1).Return(beer, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/beers/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var respBeer beers.Beer
		json.Unmarshal(w.Body.Bytes(), &respBeer)
		assert.Equal(t, *beer, respBeer)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/beers/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("FindBeerByID", mock.Anything, 2).Return(nil, beers.NewDomainError("BEER_NOT_FOUND", "not found", nil)).Once()

		req, _ := http.NewRequest(http.MethodGet, "/beers/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllBeers(t *testing.T) {
	mockService := new(MockBeerService)
	log := logger.NewNoOpLogger()
	handler := NewBeerHandler(mockService, log)

	r := setupRouter()
	r.GET(beersEndpoint, handler.GetAllBeers)

	t.Run("success", func(t *testing.T) {
		beersList := []beers.Beer{{ID: 1, Name: testBeerName}}
		mockService.On("FindAllBeers", mock.Anything).Return(beersList, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, beersEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var respBeers []beers.Beer
		json.Unmarshal(w.Body.Bytes(), &respBeers)
		assert.Equal(t, beersList, respBeers)
		mockService.AssertExpectations(t)
	})

	t.Run(serviceErr, func(t *testing.T) {
		mockService.On("FindAllBeers", mock.Anything).Return([]beers.Beer{}, errors.New(serviceErr)).Once()

		req, _ := http.NewRequest(http.MethodGet, beersEndpoint, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCalculateBoxPrice(t *testing.T) {
	mockService := new(MockBeerService)
	log := logger.NewNoOpLogger()
	handler := NewBeerHandler(mockService, log)

	r := setupRouter()
	r.GET("/beers/:id/boxprice", handler.CalculateBoxPrice)

	t.Run("success", func(t *testing.T) {
		boxPrice := &primary.BoxPriceResponse{TotalPrice: 24.0}
		mockService.On("CalculateBoxPrice", mock.Anything, mock.Anything).Return(boxPrice, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/beers/1/boxprice?quantity=6&currency=USD", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var respBoxPrice primary.BoxPriceResponse
		json.Unmarshal(w.Body.Bytes(), &respBoxPrice)
		assert.Equal(t, *boxPrice, respBoxPrice)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid quantity", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/beers/1/boxprice?quantity=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
