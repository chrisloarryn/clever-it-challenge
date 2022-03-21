package main

import (
	"fmt"

	"CleverIT-challenge/internal/http/server"
	"CleverIT-challenge/internal/infrastructure/dependencies"
)

func main() {
	fmt.Println("Running...")

	container := dependencies.NewContainer()

	server.Run(container)
}
