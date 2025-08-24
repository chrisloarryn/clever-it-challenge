package dependencies

import (
	"context"
	"fmt"

	httpAdapter "beers-challenge/internal/adapters/http"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/core/ports/secondary"
	"beers-challenge/internal/core/services"
	"beers-challenge/internal/infrastructure/config"
	"beers-challenge/internal/infrastructure/external/currencyLayer"
	"beers-challenge/internal/infrastructure/logger"
	"beers-challenge/internal/infrastructure/storage"
)

// Container holds all application dependencies
type Container struct {
	// Configuration
	config *config.ConfigProvider

	// Infrastructure
	logger          secondary.Logger
	beerRepository  secondary.BeerRepository
	currencyService secondary.CurrencyService

	// Services
	beerService primary.BeerService

	// Adapters
	httpServer *httpAdapter.Server
}

// NewContainer creates and configures a new dependency injection container
func NewContainer() (*Container, error) {
	container := &Container{}

	// Initialize configuration
	if err := container.initConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Initialize infrastructure
	if err := container.initInfrastructure(); err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	// Initialize services
	if err := container.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize adapters
	if err := container.initAdapters(); err != nil {
		return nil, fmt.Errorf("failed to initialize adapters: %w", err)
	}

	return container, nil
}

// initConfig initializes the configuration
func (c *Container) initConfig() error {
	c.config = config.NewConfigProvider()
	return nil
}

// initInfrastructure initializes infrastructure components
func (c *Container) initInfrastructure() error {
	// Initialize logger
	c.logger = logger.NewStructuredLogger(
		c.config.GetString("logger.level"),
		c.config.GetString("logger.format"),
	)

	// Initialize repository
	repositoryFactory := storage.NewRepositoryFactory(c.config)
	var err error
	c.beerRepository, err = repositoryFactory.CreateBeerRepository()
	if err != nil {
		return fmt.Errorf("failed to create beer repository: %w", err)
	}

	// Initialize currency service
	c.currencyService = currencyLayer.NewCurrencyService(c.config)

	return nil
}

// initServices initializes business services
func (c *Container) initServices() error {
	c.beerService = services.NewBeerService(
		c.beerRepository,
		c.currencyService,
		c.logger,
	)

	return nil
}

// initAdapters initializes adapters
func (c *Container) initAdapters() error {
	c.httpServer = httpAdapter.NewServer(
		c.beerService,
		c.config,
		c.logger,
	)

	return nil
}

// GetHTTPServer returns the HTTP server
func (c *Container) GetHTTPServer() *httpAdapter.Server {
	return c.httpServer
}

// GetLogger returns the logger
func (c *Container) GetLogger() secondary.Logger {
	return c.logger
}

// GetConfig returns the configuration provider
func (c *Container) GetConfig() *config.ConfigProvider {
	return c.config
}

// GetBeerService returns the beer service
func (c *Container) GetBeerService() primary.BeerService {
	return c.beerService
}

// GetBeerRepository returns the beer repository
func (c *Container) GetBeerRepository() secondary.BeerRepository {
	return c.beerRepository
}

// GetCurrencyService returns the currency service
func (c *Container) GetCurrencyService() secondary.CurrencyService {
	return c.currencyService
}

// Close gracefully closes all resources
func (c *Container) Close() error {
	ctx := context.TODO()
	c.logger.Info(ctx, "Closing container resources", nil)

	// Close repository if it has a Close method
	if closer, ok := c.beerRepository.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			c.logger.Error(ctx, "Failed to close repository", err, nil)
			return err
		}
	}

	return nil
}
