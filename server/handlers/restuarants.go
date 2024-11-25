package handlers

import (
	"context"
	"net/http"
	"time"

	"server/db"
	"server/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var restaurantCollection = db.GetCollection("taskManager", "restaurants")

func GetRestaurants(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := restaurantCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching restaurants"})
		return
	}

	var restaurants []models.Restaurant
	if err := cursor.All(ctx, &restaurants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing restaurants"})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

func GetRestaurantByID(c *gin.Context) {
	id := c.Param("id")
	objID, _ := primitive.ObjectIDFromHex(id)

	var restaurant models.Restaurant
	err := restaurantCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&restaurant)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := restaurantCollection.InsertOne(context.Background(), restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
