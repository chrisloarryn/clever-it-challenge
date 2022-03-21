package main

import (
	"fmt"

	"CleverIT-challenge/internal/http/client/currencyLayer"
	"CleverIT-challenge/internal/http/server"
	"CleverIT-challenge/internal/infrastructure/dependencies"
	"CleverIT-challenge/internal/infrastructure/storage/inmemory"
)

func main() {
	fmt.Println("Running...")
	service := currencyLayer.NewCurrencyService()

	container := dependencies.Container{
		CurrencyService: service,
		BeersRepository: inmemory.NewInMemoryRepository(),
	}

	server.Run(container)
}
