package models

// Restaurant represents the structure of a restaurant in the database and API
type Restaurant struct {
	ID           string            `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string            `json:"name" bson:"name"`
	Address      string            `json:"address" bson:"address"`
	Phone        string            `json:"phone" bson:"phone"`
	Website      string            `json:"website" bson:"website"`
	OpeningHours map[string]string `json:"opening_hours" bson:"opening_hours"`
	CuisineType  string            `json:"cuisine_type" bson:"cuisine_type"`
	IsKosher     bool              `json:"is_kosher" bson:"is_kosher"`
}
