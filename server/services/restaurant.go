package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"server/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SearchFilters struct {
	Cuisine  string
	IsKosher string
	IsOpen   string
}

func SearchRestaurants(ctx context.Context, client *dynamodb.Client, tableName string, filters SearchFilters) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant

	// Initialize ScanInput
	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	// Handle pagination
	for {
		result, err := client.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var batch []models.Restaurant
		err = attributevalue.UnmarshalListOfMaps(result.Items, &batch)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, batch...)

		if result.LastEvaluatedKey == nil {
			break
		}
		input.ExclusiveStartKey = result.LastEvaluatedKey
	}

	// Apply in-memory filtering
	filtered := filterRestaurants(restaurants, filters)

	if len(filtered) == 0 {
		return nil, errors.New("no matching restaurants found")
	}

	return filtered, nil
}

func filterRestaurants(restaurants []models.Restaurant, filters SearchFilters) []models.Restaurant {
	var filtered []models.Restaurant
	for _, r := range restaurants {
		// Filter by Cuisine
		if filters.Cuisine != "" && !strings.EqualFold(r.CuisineType, filters.Cuisine) {
			continue
		}

		// Filter by Kosher
		if filters.IsKosher != "" {
			isKosher := strings.EqualFold(filters.IsKosher, "true")
			if r.IsKosher != isKosher {
				continue
			}
		}

		// Filter by Currently Open
		if filters.IsOpen != "" && strings.EqualFold(filters.IsOpen, "true") {
			if !isRestaurantOpen(r.OpeningHours) {
				continue
			}
		}

		// If all filters match, add the restaurant
		filtered = append(filtered, r)
	}

	return filtered
}

func isRestaurantOpen(openingHours map[string]string) bool {
	// Get current day and time
	now := time.Now()
	currentDay := now.Weekday().String()
	currentTime := now.Format("15:04")

	// Check if the restaurant is open on the current day
	hours, ok := openingHours[currentDay]
	if !ok || hours == "Closed" {
		return false
	}

	// Parse opening and closing times
	times := strings.Split(hours, "-")
	if len(times) != 2 {
		return false
	}

	openTime, err := time.Parse("15:04", times[0])
	if err != nil {
		return false
	}
	closeTime, err := time.Parse("15:04", times[1])
	if err != nil {
		return false
	}

	return currentTime >= openTime.Format("15:04") && currentTime <= closeTime.Format("15:04")
}

func AddRestaurant(ctx context.Context, client *dynamodb.Client, tableName string, restaurant models.Restaurant) error {
	// Convert restaurant to DynamoDB item
	item, err := attributevalue.MarshalMap(restaurant)
	if err != nil {
		return err
	}

	// Add the item to DynamoDB
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return err
}

func RemoveRestaurant(ctx context.Context, client *dynamodb.Client, tableName string, restaurantID string) error {
	// Define the key for the item to delete
	key := map[string]types.AttributeValue{
		"restaurant_id": &types.AttributeValueMemberS{Value: restaurantID},
	}

	// Remove the item from DynamoDB
	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       key,
	})
	return err
}

func EditRestaurant(ctx context.Context, client *dynamodb.Client, tableName string, restaurant models.Restaurant) error {
	// Convert updated restaurant to DynamoDB item
	item, err := attributevalue.MarshalMap(restaurant)
	if err != nil {
		return err
	}

	// Replace the existing item in DynamoDB
	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	return err
}
