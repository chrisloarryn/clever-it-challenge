package server

import (
	"CleverIT-challenge/internal/http/server/handlers"
	"CleverIT-challenge/internal/infrastructure/dependencies"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {

}

func Run(container dependencies.Container) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	findAllBeersHandler := handlers.NewFindAllBeersHandler(container)
	r.GET("/beers", findAllBeersHandler.GetAllBeers)
	createBeerHandler := handlers.NewCreateBeerHandler(container)
	r.POST("/beers", createBeerHandler.CreateBeer)
	findOneBeerHandler := handlers.NewFindOneBeerHandler(container)
	r.GET("/beers/:beerID", findOneBeerHandler.FindOneBeer)
	calculateBeerBoxHandler := handlers.NewCalculateBeerBoxHandler(container)
	r.GET("/beers/:beerID/boxprice", calculateBeerBoxHandler.CalculateBeerBox)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
