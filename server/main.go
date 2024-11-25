package main

import (
	"context"
	"log"
	"os"

	"server/data"
	"server/handlers"
	"server/middleware"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	svc       *dynamodb.Client
	tableName = "restaurants" // Move to an env variable if needed
)

func init() {
	// Initialize the DynamoDB client
	initializeDynamoDB()

	// Populate the table if it is empty
	populateTableIfEmpty()
}

func initializeDynamoDB() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}
	svc = dynamodb.NewFromConfig(cfg)
	log.Println("Successfully initialized DynamoDB client")
}

func populateTableIfEmpty() {
	// Check if the table is populated
	populated, err := data.IsTablePopulated(context.TODO(), svc, tableName)
	if err != nil {
		log.Fatalf("Error checking table population: %v", err)
	}
	if populated {
		log.Printf("Table %s is already populated. Skipping initialization.", tableName)
		return
	}

	log.Printf("Table %s is empty. Initializing with data.", tableName)

	// Load restaurant data from JSON file
	restaurants, err := data.LoadRestaurants("data/restuarants_data.json")
	if err != nil {
		log.Fatalf("Failed to load restaurants from JSON: %v", err)
	}

	// Insert restaurant data into DynamoDB table
	err = data.InsertRestaurants(context.TODO(), svc, tableName, restaurants)
	if err != nil {
		log.Fatalf("Failed to insert restaurants: %v", err)
	}

	log.Println("Successfully populated DynamoDB table with restaurant data")
}

func setupRoutes(client *dynamodb.Client) *gin.Engine {
	r := gin.Default()

	// Middleware for logging searches
	r.Use(middleware.AuditLog(client))

	// Public routes
	r.GET("/restaurants/search", func(c *gin.Context) {
		handlers.SearchRestaurants(c, client)
	})

	// Admin routes
	admin := r.Group("/admin", handlers.AdminAuthMiddleware())
	{
		admin.POST("/restaurants", func(c *gin.Context) {
			handlers.AddRestaurant(c, client)
		})
		admin.PUT("/restaurants", func(c *gin.Context) {
			handlers.EditRestaurant(c, client)
		})
		admin.DELETE("/restaurants/:id", func(c *gin.Context) {
			handlers.RemoveRestaurant(c, client)
		})
		admin.GET("/audit-logs", func(c *gin.Context) {
			handlers.FetchAuditLogs(c, client)
		})
	}

	return r
}

func main() {
	// Initialize Gin routes with the DynamoDB client
	r := setupRoutes(svc)

	r.StaticFile("/admin", "./fe/admin.html")

	// Determine the port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Start the server
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
