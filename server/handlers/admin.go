package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"server/models"
	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func generateUniqueID() string {
	return uuid.New().String()
}

func AddRestaurant(c *gin.Context, client *dynamodb.Client) {
	var restaurant models.Restaurant

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant data", "details": err.Error()})
		return
	}

	if restaurant.RestaurantID == "" {
		restaurant.RestaurantID = generateUniqueID()
	}

	log.Printf("Restaurant to be added: %+v", restaurant)

	// Call the service to add the restaurant
	err := services.AddRestaurant(c.Request.Context(), client, "restaurants", restaurant)
	if err != nil {
		log.Printf("Error inserting restaurant: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add restaurant", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant added successfully"})
}

func RemoveRestaurant(c *gin.Context, client *dynamodb.Client) {
	restaurantID := c.Param("id")

	// Remove the restaurant from DynamoDB
	err := services.RemoveRestaurant(c.Request.Context(), client, "restaurants", restaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant removed successfully"})
}

func EditRestaurant(c *gin.Context, client *dynamodb.Client) {
	restaurantID := c.Param("id")
	var restaurant models.Restaurant

	// Bind JSON payload to restaurant struct
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant data"})
		return
	}

	restaurant.RestaurantID = restaurantID // Ensure the correct restaurant_id is set

	// Update the restaurant in DynamoDB
	err := services.EditRestaurant(c.Request.Context(), client, "restaurants", restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
}

func GetRestaurantByID(c *gin.Context, client *dynamodb.Client) {
	restaurantID := c.Param("id")
	restaurant, err := services.FetchRestaurantByID(c.Request.Context(), client, "restaurants", restaurantID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurant details"})
		return
	}

	if restaurant == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// AdminAuthMiddleware protects admin routes with a password
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the password from the Authorization header
		providedPassword := c.GetHeader("Authorization")

		// Get the expected admin password
		expectedPassword := os.Getenv("ADMIN_PASSWORD")
		log.Printf("Expected admin password: %s", expectedPassword)
		if expectedPassword == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server is not configured properly"})
			return
		}

		// Check if the provided password matches
		if providedPassword != expectedPassword {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Continue to the next handler if authorized
		c.Next()
	}
}

func GetAuditLogsHandler(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the 'minutes' query parameter (optional)
		minutesParam := c.Query("minutes")
		minutes := 0 // Default to 0 (which means last 24 hours in the GetFilteredLogs function)

		if minutesParam != "" {
			var err error
			minutes, err = strconv.Atoi(minutesParam)
			if err != nil || minutes < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'minutes' parameter. It must be a positive integer."})
				return
			}
		}

		// Call the service function to fetch logs
		logs, err := services.GetFilteredLogs(c.Request.Context(), client, minutes) // Correctly call the function from the `services` package
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
			return
		}

		// Return the logs as JSON
		c.JSON(http.StatusOK, logs)
	}
}
