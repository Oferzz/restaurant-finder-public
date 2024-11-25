package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	log.Println("Connected to MongoDB")
	Client = client
}

func GetCollection(database, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}
