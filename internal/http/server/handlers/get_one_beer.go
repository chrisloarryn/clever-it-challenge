package handlers

import (
	"CleverIT-challenge/internal/core/usecases"
	"CleverIT-challenge/internal/infrastructure/dependencies"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FindOneBeerHandler struct {
	uc *usecases.FinderBeersByID
}

func NewFindOneBeerHandler(container dependencies.Container) *FindOneBeerHandler {
	return &FindOneBeerHandler{
		uc: usecases.NewFinderBeersByID(container.BeersRepository),
	}
}

func (handler *FindOneBeerHandler) FindOneBeer(ctx *gin.Context) {
	beerID, err := strconv.Atoi(ctx.Param("beerID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Invalid beer ID: %s", err.Error()))
		return
	}
	beer, err := handler.uc.Execute(ctx, beerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, beer)

}
