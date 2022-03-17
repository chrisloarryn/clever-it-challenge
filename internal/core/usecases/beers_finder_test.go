package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/beers/beersmocks"
	"CleverIT-challenge/internal/core/usecases"
)

func TestFinderAllBeers_Execute_ShouldReturnsABeerList(t *testing.T) {
	t.Log("Should returns a beer list")
	// Setup
	controller := gomock.NewController(t)

	repository := beersmocks.NewMockRepository(controller)
	beersList := []beers.Beer{
		{
			ID:       123,
			Name:     "Golden",
			Brewery:  "Kross",
			Price:    10.5,
			Currency: "EUR",
			Country:  "Chile",
		},
		{},
	}
	repository.EXPECT().FindAllBeers(gomock.Any()).Return(beersList, nil).Times(1)

	finderAllBeers := usecases.NewFinderAllBeers(repository)

	// Execute
	result, err := finderAllBeers.Execute(context.TODO())

	// Verify
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, beersList, result)
}

func TestFinderAllBeers_Execute_ShouldReturnsAnErrorInRepository(t *testing.T) {
	t.Log("Should returns a beer list")
	// Setup
	controller := gomock.NewController(t)

	repository := beersmocks.NewMockRepository(controller)
	customError := fmt.Errorf("this is a custom error")
	repository.EXPECT().FindAllBeers(gomock.Any()).Return(nil, customError).Times(1)

	finderAllBeers := usecases.NewFinderAllBeers(repository)

	// Execute
	result, err := finderAllBeers.Execute(context.TODO())

	// Verify
	require.Nil(t, result)
	assert.EqualError(t, err, customError.Error())
}
