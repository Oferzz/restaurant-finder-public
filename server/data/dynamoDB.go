package data

import (
	"context"
	"log"

	"server/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func IsTablePopulated(ctx context.Context, svc *dynamodb.Client, tableName string) (bool, error) {
	input := &dynamodb.ScanInput{
		TableName: &tableName,
		Limit:     aws.Int32(1), // Corrected
		Select:    types.SelectCount,
	}

	result, err := svc.Scan(ctx, input)
	if err != nil {
		log.Printf("Error scanning table: %v", err)
		return false, err
	}

	return result.Count > 0, nil // Corrected
}

func InsertRestaurants(ctx context.Context, svc *dynamodb.Client, tableName string, restaurants []models.Restaurant) error {
	for _, restaurant := range restaurants {
		_, err := svc.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]types.AttributeValue{
				"restaurant_id":   &types.AttributeValueMemberS{Value: restaurant.RestaurantID},
				"restaurant_name": &types.AttributeValueMemberS{Value: restaurant.Name},
				"address":         &types.AttributeValueMemberS{Value: restaurant.Address},
				"phone":           &types.AttributeValueMemberS{Value: restaurant.Phone},
				"website":         &types.AttributeValueMemberS{Value: restaurant.Website},
				"cuisine_type":    &types.AttributeValueMemberS{Value: restaurant.CuisineType},
				"is_kosher":       &types.AttributeValueMemberBOOL{Value: restaurant.IsKosher},
				"opening_hours": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
					"Monday":    &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Monday"]},
					"Tuesday":   &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Tuesday"]},
					"Wednesday": &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Wednesday"]},
					"Thursday":  &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Thursday"]},
					"Friday":    &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Friday"]},
					"Saturday":  &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Saturday"]},
					"Sunday":    &types.AttributeValueMemberS{Value: restaurant.OpeningHours["Sunday"]},
				}},
			},
		})
		if err != nil {
			log.Printf("Failed to insert restaurant %s: %v", restaurant.Name, err)
			return err
		}
		log.Printf("Inserted restaurant: %s", restaurant.Name)
	}
	return nil
}
