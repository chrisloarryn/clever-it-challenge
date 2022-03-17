package usecases_test

import (
	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/beers/beersmocks"
	"CleverIT-challenge/internal/core/usecases"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFinderBeersByID_Execute_ShouldReturnsABeerData(t *testing.T) {
	t.Log("Should returns a beer from his ID")
	// Setup
	controller := gomock.NewController(t)

	beerID := 123
	beerResult := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    10.5,
		Currency: "EUR",
		Country:  "Chile",
	}

	repository := beersmocks.NewMockRepository(controller)
	repository.EXPECT().FindBeerByID(gomock.Any(), beerID).Return(beerResult, nil).Times(1)

	finderAllBeers := usecases.NewFinderBeersByID(repository)

	// Execute
	result, err := finderAllBeers.Execute(context.TODO(), beerID)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, beerResult, result)
}

func TestFinderBeersByID_Execute_ShouldReturnsAnErrorFromRepository(t *testing.T) {
	t.Log("Should returns an error from repository")
	// Setup
	controller := gomock.NewController(t)

	beerID := 123
	customError := fmt.Errorf("this is a custom error")

	repository := beersmocks.NewMockRepository(controller)
	repository.EXPECT().FindBeerByID(gomock.Any(), beerID).Return(beers.Beer{}, customError).Times(1)

	finderAllBeers := usecases.NewFinderBeersByID(repository)

	// Execute
	result, err := finderAllBeers.Execute(context.TODO(), beerID)

	// Verify
	require.Error(t, err, customError.Error())
	assert.Equal(t, beers.Beer{},result)
}

func TestFinderBeersByID_Execute_ShouldReturnsAnErrorForInvalidID(t *testing.T) {
	t.Log("Should returns an error for invalid ID")
	// Setup
	controller := gomock.NewController(t)

	beerID := -1
	repository := beersmocks.NewMockRepository(controller)

	finderAllBeers := usecases.NewFinderBeersByID(repository)

	// Execute
	result, err := finderAllBeers.Execute(context.TODO(), beerID)

	// Verify
	require.Error(t, err, "Invalid ID: -1")
	assert.Equal(t, beers.Beer{},result)
}
