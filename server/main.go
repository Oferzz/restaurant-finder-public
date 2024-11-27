package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"server/data"
	"server/handlers"
	"server/middleware"
	"server/services"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	svc           *dynamodb.Client
	tableName     = "restaurants"
	adminPassword string
)

func init() {
	// Load environment variables
	loadEnvironmentVariables()

	// Initialize the DynamoDB client
	initializeDynamoDB()

	// Populate the table if it is empty
	populateTableIfEmpty()
}

func loadEnvironmentVariables() {
	// Load admin password from environment variable
	adminPassword = os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatalf("ADMIN_PASSWORD environment variable is not set")
	}

	// Optionally override the table name from an environment variable
	if envTableName := os.Getenv("TABLE_NAME"); envTableName != "" {
		tableName = envTableName
	}
	log.Printf("Using table name: %s", tableName)
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
	restaurants, err := data.LoadRestaurants("data/restaurants_data.json")
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

	// Add middleware
	r.Use(gin.Logger())                // Request logging
	r.Use(middleware.AuditLog(client)) // Audit logging for searches

	r.GET("/readiness", func(c *gin.Context) {
		log.Println("Readiness check triggered")
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/liveness", func(c *gin.Context) {
		log.Println("Liveness check triggered")
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Public routes
	setupPublicRoutes(r, client)

	// Admin routes
	setupAdminRoutes(r, client)

	return r
}

func setupPublicRoutes(r *gin.Engine, client *dynamodb.Client) {
	r.GET("/restaurants/search", func(c *gin.Context) {
		handlers.SearchRestaurants(c, client)
	})
}

func setupAdminRoutes(r *gin.Engine, client *dynamodb.Client) {
	admin := r.Group("/admin", handlers.AdminAuthMiddleware())
	{
		admin.GET("/validate", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Password is valid"})
		})
		admin.POST("/restaurants", func(c *gin.Context) {
			handlers.AddRestaurant(c, client)
		})
		admin.PUT("/restaurants/:id", func(c *gin.Context) {
			handlers.EditRestaurant(c, client)
		})
		admin.DELETE("/restaurants/:id", func(c *gin.Context) {
			handlers.RemoveRestaurant(c, client)
		})
		admin.GET("/logs", func(c *gin.Context) {
			// Fetch query parameter for 'minutes'
			minutesParam := c.DefaultQuery("minutes", "1440") // Default to 1440 minutes (24 hours)
			minutes, err := strconv.Atoi(minutesParam)
			if err != nil || minutes < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minutes parameter"})
				return
			}

			// Use the appropriate function from the services package
			logs, err := services.GetFilteredLogs(c.Request.Context(), client, minutes)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
				return
			}

			c.JSON(http.StatusOK, logs)
		})
		admin.GET("/restaurants/:id", func(c *gin.Context) {
			handlers.GetRestaurantByID(c, client)
		})
	}
}

func main() {
	// Initialize Gin routes with the DynamoDB client
	r := setupRoutes(svc)

	// Static file serving
	r.Static("/static", "./static")
	r.StaticFile("/admin", "./static/admin.html")

	// Determine the port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown the server
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
