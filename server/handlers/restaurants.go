package handlers

import (
	"log"
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

	// Validate query parameters
	if isKosher != "" && isKosher != "true" && isKosher != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'is_kosher'. Must be 'true' or 'false'."})
		return
	}
	if isOpen != "" && isOpen != "true" && isOpen != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'is_open'. Must be 'true' or 'false'."})
		return
	}

	// Create filters
	filters := services.SearchFilters{
		Cuisine:  cuisine,
		IsKosher: isKosher,
		IsOpen:   isOpen,
	}

	// Call service function
	restaurants, err := services.SearchRestaurants(c.Request.Context(), client, "restaurants", filters)
	if err != nil {
		log.Printf("Error searching restaurants: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants. Please try again later."})
		return
	}

	// Handle empty results
	if len(restaurants) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No restaurants match the given criteria."})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, restaurants)
}
