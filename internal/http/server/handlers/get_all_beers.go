package handlers

import (
	"CleverIT-challenge/internal/core/usecases"
	"CleverIT-challenge/internal/infrastructure/dependencies"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FindAllBeersHandler struct {
	uc *usecases.FinderAllBeers
}

func NewFindAllBeersHandler(container dependencies.Container) *FindAllBeersHandler {
	return &FindAllBeersHandler{
		uc: usecases.NewFinderAllBeers(container.BeersRepository()),
	}
}

func (handler *FindAllBeersHandler) GetAllBeers(ctx *gin.Context) {
	beers, err := handler.uc.Execute(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, beers)

}
