package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"CleverIT-challenge/internal/core/usecases"
	"CleverIT-challenge/internal/infrastructure/dependencies"
)

type CalculateBeerBoxHandler struct {
	uc *usecases.BoxPriceCalculator
}

func NewCalculateBeerBoxHandler(container dependencies.Container) *CalculateBeerBoxHandler {
	return &CalculateBeerBoxHandler{
		uc: usecases.NewBoxPriceCalculator(container.BeersRepository, container.CurrencyService),
	}
}

func (handler *CalculateBeerBoxHandler) CalculateBeerBox(ctx *gin.Context) {
	beerID, err := strconv.Atoi(ctx.Param("beerID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid beer ID: %s", err.Error()))
		return
	}
	currency := ctx.Query("currency")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid beer ID: %s", err.Error()))
		return
	}
	quantity, err := strconv.Atoi(ctx.DefaultQuery("quantity", "6"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid beer ID: %s", err.Error()))
		return
	}

	price, err := handler.uc.Execute(ctx, beerID, quantity, currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Amount float64 `json:"amount"`
	}{
		Amount: price,
	})
}
