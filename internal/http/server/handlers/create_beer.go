package handlers

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/usecases"
	"CleverIT-challenge/internal/infrastructure/dependencies"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateBeerHandler struct {
	uc *usecases.CreateBeer
}

func NewCreateBeerHandler(container dependencies.Container) *CreateBeerHandler {
	return &CreateBeerHandler{
		uc: usecases.NewCreateBeer(container.BeersRepository(), container.CurrencyService()),
	}
}

func (handler *CreateBeerHandler) CreateBeer(ctx *gin.Context) {
	beer := beers.Beer{}
	if err := ctx.BindJSON(&beer); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := handler.uc.Execute(ctx, beer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Status(http.StatusCreated)

}
