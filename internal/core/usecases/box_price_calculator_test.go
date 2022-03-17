package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"CleverIT-challenge/internal/core/domain/beers"
	"CleverIT-challenge/internal/core/domain/beers/beersmocks"
	"CleverIT-challenge/internal/core/domain/currency/currencymocks"
	"CleverIT-challenge/internal/core/usecases"
)

func TestBoxPriceCalculator_Execute_ShouldReturnAPriceFor10BoxesFromUSDtoUSD(t *testing.T) {
	t.Log("Should returns a price for 10 boxes from USD to USD")

	controller := gomock.NewController(t)
	ctx := context.TODO()
	beerID := 123
	beer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    10.5,
		Currency: "USD",
		Country:  "Chile",
	}
	quantity := 10
	currencyTo := "USD"

	currencyService := currencymocks.NewMockService(controller)
	beersRepository := beersmocks.NewMockRepository(controller)

	beersRepository.EXPECT().FindBeerByID(ctx, beerID).Return(beer, nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, beer.Currency).Return(float64(1), nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, currencyTo).Return(float64(1), nil)

	boxCalculator := usecases.NewBoxPriceCalculator(beersRepository, currencyService)

	boxPrice, err := boxCalculator.Execute(ctx, beerID, quantity, currencyTo)

	require.NoError(t, err)
	assert.Equal(t, boxPrice, float64(10))
}

func TestBoxPriceCalculator_Execute_ShouldReturnAPriceFor10BoxesFromARStoUSD(t *testing.T) {
	t.Log("Should returns a price for 10 boxes from ARS to USD")

	controller := gomock.NewController(t)
	ctx := context.TODO()
	beerID := 123
	beer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    200,
		Currency: "ARS",
		Country:  "Chile",
	}
	quantity := 10
	currencyTo := "USD"

	currencyService := currencymocks.NewMockService(controller)
	beersRepository := beersmocks.NewMockRepository(controller)

	beersRepository.EXPECT().FindBeerByID(ctx, beerID).Return(beer, nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, beer.Currency).Return(float64(109), nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, currencyTo).Return(float64(1), nil)

	boxCalculator := usecases.NewBoxPriceCalculator(beersRepository, currencyService)

	boxPrice, err := boxCalculator.Execute(ctx, beerID, quantity, currencyTo)

	require.NoError(t, err)
	assert.Equal(t, boxPrice, 0.09174311926605505)
}

func TestBoxPriceCalculator_Execute_ShouldReturnAPriceFor10BoxesFromARStoEUR(t *testing.T) {
	t.Log("Should returns a price for 10 boxes from ARS to EUR")

	controller := gomock.NewController(t)
	ctx := context.TODO()
	beerID := 123
	beer := beers.Beer{
		ID:       123,
		Name:     "Golden",
		Brewery:  "Kross",
		Price:    200,
		Currency: "ARS",
		Country:  "Chile",
	}
	quantity := 10
	currencyTo := "EUR"

	currencyService := currencymocks.NewMockService(controller)
	beersRepository := beersmocks.NewMockRepository(controller)

	beersRepository.EXPECT().FindBeerByID(ctx, beerID).Return(beer, nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, beer.Currency).Return(float64(109), nil)
	currencyService.EXPECT().GetCurrencyPriceInDollar(ctx, currencyTo).Return(0.9, nil)

	boxCalculator := usecases.NewBoxPriceCalculator(beersRepository, currencyService)

	boxPrice, err := boxCalculator.Execute(ctx, beerID, quantity, currencyTo)

	require.NoError(t, err)
	assert.Equal(t, boxPrice, 0.08256880733944955)
}

func TestBoxPriceCalculator_Execute_ShouldFailForNegativeQuantity(t *testing.T) {
	t.Log("Should fail for negative quantity")

	controller := gomock.NewController(t)
	ctx := context.TODO()
	beerID := 123
	quantity := -1
	currencyTo := "EUR"

	currencyService := currencymocks.NewMockService(controller)
	beersRepository := beersmocks.NewMockRepository(controller)

	boxCalculator := usecases.NewBoxPriceCalculator(beersRepository, currencyService)

	boxPrice, err := boxCalculator.Execute(ctx, beerID, quantity, currencyTo)

	require.Errorf(t, err, "invalid quantity")
	assert.Equal(t, boxPrice, float64(0))
}
