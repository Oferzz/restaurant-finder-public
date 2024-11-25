package handlers

import (
	"net/http"

	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func SearchRestaurants(c *gin.Context, client *dynamodb.Client) {
	// Get query parameters
	cuisine := c.Query("cuisine")
	isKosher := c.Query("is_kosher")
	isOpen := c.Query("is_open")

	// Create filters
	filters := services.SearchFilters{
		Cuisine:  cuisine,
		IsKosher: isKosher,
		IsOpen:   isOpen,
	}

	// Call service function
	restaurants, err := services.SearchRestaurants(c.Request.Context(), client, "restaurants", filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}
