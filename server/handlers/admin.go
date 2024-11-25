package handlers

import (
	"net/http"

	"server/models"
	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AddRestaurant(c *gin.Context, client *dynamodb.Client) {
	var restaurant models.Restaurant
	if err := c.BindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant data"})
		return
	}

	err := services.AddRestaurant(c.Request.Context(), client, "restaurants", restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant added successfully"})
}

func RemoveRestaurant(c *gin.Context, client *dynamodb.Client) {
	id := c.Param("id")

	err := services.RemoveRestaurant(c.Request.Context(), client, "restaurants", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant removed successfully"})
}

func EditRestaurant(c *gin.Context, client *dynamodb.Client) {
	var restaurant models.Restaurant
	if err := c.BindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant data"})
		return
	}

	err := services.EditRestaurant(c.Request.Context(), client, "restaurants", restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit restaurant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
}

// AdminAuthMiddleware protects admin routes with a password
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminPassword := "admin"
		if adminPassword == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin password not configured"})
			c.Abort()
			return
		}

		// Get the password from query or header
		password := c.Query("password")
		if password == "" {
			password = c.GetHeader("Authorization")
		}

		// Validate the password
		if password != adminPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort()
			return
		}

		// Continue to the next handler if authenticated
		c.Next()
	}
}

func FetchAuditLogs(c *gin.Context, client *dynamodb.Client) {
	// Retrieve audit logs from services
	logs, err := services.GetAuditLogs(c.Request.Context(), client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
