package models

type Restaurant struct {
	RestaurantID string            `json:"restaurant_id" dynamodbav:"restaurant_id"`
	Name         string            `json:"restaurant_name" dynamodbav:"restaurant_name"`
	Address      string            `json:"address" dynamodbav:"address"`
	Phone        string            `json:"phone" dynamodbav:"phone"`
	Website      string            `json:"website" dynamodbav:"website"`
	CuisineType  string            `json:"cuisine_type" dynamodbav:"cuisine_type"`
	IsKosher     bool              `json:"is_kosher" dynamodbav:"is_kosher"`
	OpeningHours map[string]string `json:"opening_hours" dynamodbav:"opening_hours"`
}
