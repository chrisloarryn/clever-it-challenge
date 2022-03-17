package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/beers/beersmocks"
	"CleverIT-challenge/internal/core/usecases"
)

func TestCreateBeer_Execute_ShouldCreateABeer(t *testing.T) {
	t.Log("Should create a Beer")
	// Setup
	controller := gomock.NewController(t)

	repository := beersmocks.NewMockRepository(controller)
	newBeer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    10.5,
		Currency: "EUR",
		Country:  "Chile",
	}
	repository.EXPECT().SaveBeer(gomock.Any(), newBeer).Return(nil).Times(1)

	createBeerUseCase := usecases.NewCreateBeer(repository)

	// Execute
	err := createBeerUseCase.Execute(context.TODO(), newBeer)

	// Verify
	assert.NoError(t, err)
}

func TestCreateBeer_Execute_ShouldReturnAnError(t *testing.T) {
	t.Log("Should return an error when try to create a beer")
	// Setup
	controller := gomock.NewController(t)
	newBeer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    10.5,
		Currency: "EUR",
		Country:  "Chile",
	}
	customError := fmt.Errorf("this is a custom error")

	repository := beersmocks.NewMockRepository(controller)
	repository.EXPECT().SaveBeer(gomock.Any(), newBeer).Return(customError).Times(1)

	createBeerUseCase := usecases.NewCreateBeer(repository)

	// Execute
	err := createBeerUseCase.Execute(context.TODO(), newBeer)

	// Verify
	assert.EqualError(t, err, customError.Error())
}

func TestCreateBeer_Execute_ShouldReturnAnErrorForInvalidNegative(t *testing.T) {
	t.Log("Should return an error when try to create a beer")
	// Setup
	controller := gomock.NewController(t)
	newBeer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    -10.5,
		Currency: "EUR",
		Country:  "Chile",
	}
	invalidPriceError := fmt.Errorf("invalid price")

	repository := beersmocks.NewMockRepository(controller)

	createBeerUseCase := usecases.NewCreateBeer(repository)

	// Execute
	err := createBeerUseCase.Execute(context.TODO(), newBeer)

	// Verify
	assert.EqualError(t, err, invalidPriceError.Error())
}
