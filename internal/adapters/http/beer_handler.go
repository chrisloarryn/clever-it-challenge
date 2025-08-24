package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/primary"
	"beers-challenge/internal/core/ports/secondary"
)

// BeerHandler handles HTTP requests for beer operations
type BeerHandler struct {
	beerService primary.BeerService
	logger      secondary.Logger
}

// NewBeerHandler creates a new beer handler
func NewBeerHandler(beerService primary.BeerService, logger secondary.Logger) *BeerHandler {
	return &BeerHandler{
		beerService: beerService,
		logger:      logger,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// CreateBeer handles POST /beers
func (h *BeerHandler) CreateBeer(c *gin.Context) {
	var req primary.CreateBeerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(c.Request.Context(), "Invalid request body", err, map[string]interface{}{
			"endpoint": "POST /beers",
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.beerService.CreateBeer(c.Request.Context(), req); err != nil {
		h.handleError(c, "Failed to create beer", err)
		return
	}

	h.logger.Info(c.Request.Context(), "Beer created successfully", map[string]interface{}{
		"beer_id": req.ID,
		"name":    req.Name,
	})

	c.Status(http.StatusCreated)
}

// GetBeer handles GET /beers/:id
func (h *BeerHandler) GetBeer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Invalid beer ID", err, map[string]interface{}{
			"id_param": idParam,
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Beer ID must be a valid integer",
		})
		return
	}

	beer, err := h.beerService.FindBeerByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, "Failed to find beer", err)
		return
	}

	c.JSON(http.StatusOK, beer)
}

// GetAllBeers handles GET /beers
func (h *BeerHandler) GetAllBeers(c *gin.Context) {
	beersSlice, err := h.beerService.FindAllBeers(c.Request.Context())
	if err != nil {
		h.handleError(c, "Failed to find beers", err)
		return
	}

	c.JSON(http.StatusOK, beersSlice)
}

// CalculateBoxPrice handles GET /beers/:id/boxprice
func (h *BeerHandler) CalculateBoxPrice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error(c.Request.Context(), "Invalid beer ID", err, map[string]interface{}{
			"id_param": idParam,
		})
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_ID",
			Message: "Beer ID must be a valid integer",
		})
		return
	}

	quantityParam := c.DefaultQuery("quantity", "1")
	quantity, err := strconv.Atoi(quantityParam)
	if err != nil || quantity < 1 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_QUANTITY",
			Message: "Quantity must be a positive integer",
		})
		return
	}

	currency := c.DefaultQuery("currency", "USD")

	req := primary.CalculateBoxPriceRequest{
		BeerID:   id,
		Quantity: quantity,
		Currency: currency,
	}

	response, err := h.beerService.CalculateBoxPrice(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, "Failed to calculate box price", err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// handleError handles errors and sends appropriate HTTP responses
func (h *BeerHandler) handleError(c *gin.Context, message string, err error) {
	h.logger.Error(c.Request.Context(), message, err, map[string]interface{}{
		"endpoint": c.Request.Method + " " + c.Request.URL.Path,
	})

	// Check if it's a validation error
	if validationErr, ok := err.(*beers.ValidationError); ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "VALIDATION_ERROR",
			Message: validationErr.Error(),
		})
		return
	}

	// Check if it's a domain error
	if domainErr, ok := err.(*beers.DomainError); ok {
		statusCode := http.StatusInternalServerError

		switch domainErr.Code {
		case "BEER_NOT_FOUND":
			statusCode = http.StatusNotFound
		case "BEER_ALREADY_EXISTS":
			statusCode = http.StatusConflict
		case "INVALID_CURRENCY":
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, ErrorResponse{
			Error:   domainErr.Code,
			Message: domainErr.Message,
		})
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	})
}
