package main

import (
	"server/db"
	"server/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to MongoDB
	db.Connect()

	r := gin.Default()

	// Define routes
	r.GET("/restaurants", handlers.GetRestaurants)
	r.GET("/restaurants/:id", handlers.GetRestaurantByID)
	r.POST("/restaurants", handlers.CreateRestaurant)

	// Start the server
	r.Run(":8080")
}
