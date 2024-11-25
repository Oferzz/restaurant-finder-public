package data

import (
	"encoding/json"
	"os"

	"server/models"
)

func LoadRestaurants(filename string) ([]models.Restaurant, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var restaurants []models.Restaurant
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&restaurants)
	if err != nil {
		return nil, err
	}

	return restaurants, nil
}
