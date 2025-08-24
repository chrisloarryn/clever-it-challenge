package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/core/ports/secondary"
	"beers-challenge/internal/infrastructure/config"
)

const (
	// API paths
	BeersPath = "/beers"
	APIPrefix = "/api/v1"
)

// Server represents the HTTP server
type Server struct {
	router      *gin.Engine
	beerHandler *BeerHandler
	config      *config.ConfigProvider
	logger      secondary.Logger
	server      *http.Server
}

// NewServer creates a new HTTP server
func NewServer(
	beerService primary.BeerService,
	config *config.ConfigProvider,
	logger secondary.Logger,
) *Server {
	// Set Gin mode based on environment
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware(logger))
	router.Use(CORSMiddleware())

	beerHandler := NewBeerHandler(beerService, logger)

	server := &Server{
		router:      router,
		beerHandler: beerHandler,
		config:      config,
		logger:      logger,
	}

	server.setupRoutes()

	return server
}

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/ping", s.healthCheck)

	// API routes
	api := s.router.Group(APIPrefix)
	{
		// Beer routes
		beers := api.Group(BeersPath)
		{
			beers.POST("", s.beerHandler.CreateBeer)
			beers.GET("", s.beerHandler.GetAllBeers)
			beers.GET("/:id", s.beerHandler.GetBeer)
			beers.GET("/:id/boxprice", s.beerHandler.CalculateBoxPrice)
		}
	}

	// Legacy routes for backward compatibility
	s.router.POST(BeersPath, s.beerHandler.CreateBeer)
	s.router.GET(BeersPath, s.beerHandler.GetAllBeers)
	s.router.GET(BeersPath+"/:id", s.beerHandler.GetBeer)
	s.router.GET(BeersPath+"/:id/boxprice", s.beerHandler.CalculateBoxPrice)
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "beer-api",
		"version":   "1.0.0",
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d",
		s.config.GetString("server.host"),
		s.config.GetInt("server.port"))

	s.server = &http.Server{
		Addr:         address,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info(context.Background(), "Starting HTTP server", map[string]interface{}{
		"address": address,
	})

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "Stopping HTTP server", nil)

	if s.server != nil {
		return s.server.Shutdown(ctx)
	}

	return nil
}

// LoggerMiddleware adds structured logging to requests
func LoggerMiddleware(logger secondary.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		fields := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		}

		if raw != "" {
			fields["query"] = raw
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		logger.Info(c.Request.Context(), "HTTP request processed", fields)
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
