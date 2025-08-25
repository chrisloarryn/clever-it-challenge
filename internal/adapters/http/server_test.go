package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/infrastructure/config"
	"beers-challenge/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBeerServiceForServer struct {
	mock.Mock
}

func (m *MockBeerServiceForServer) CreateBeer(ctx context.Context, req primary.CreateBeerRequest) error {
	return nil
}
func (m *MockBeerServiceForServer) FindBeerByID(ctx context.Context, id int) (*beers.Beer, error) {
	return nil, nil
}
func (m *MockBeerServiceForServer) FindAllBeers(ctx context.Context) ([]beers.Beer, error) {
	return nil, nil
}
func (m *MockBeerServiceForServer) CalculateBoxPrice(ctx context.Context, req primary.CalculateBoxPriceRequest) (*primary.BoxPriceResponse, error) {
	return nil, nil
}

func TestNewServer(t *testing.T) {
	cfg := config.NewConfigProvider()
	log := logger.NewNoOpLogger()
	service := new(MockBeerServiceForServer)

	server := NewServer(service, cfg, log)
	assert.NotNil(t, server)
}

func TestHealthCheck(t *testing.T) {
	cfg := config.NewConfigProvider()
	log := logger.NewNoOpLogger()
	service := new(MockBeerServiceForServer)

	server := NewServer(service, cfg, log)
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestLoggerMiddleware(t *testing.T) {
	log := logger.NewNoOpLogger()
	router := gin.New()
	router.Use(LoggerMiddleware(log))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(CORSMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestServerStop(t *testing.T) {
	cfg := config.NewConfigProvider()
	log := logger.NewNoOpLogger()
	service := new(MockBeerServiceForServer)

	server := NewServer(service, cfg, log)
	// Just call stop, we can't easily test the shutdown process here
	err := server.Stop(context.Background())
	assert.NoError(t, err)
}
