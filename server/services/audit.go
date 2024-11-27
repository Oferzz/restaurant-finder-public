package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GeoResponse struct {
	Country string `json:"country"`
}

// GetCountryFromIP fetches the country for a given IP address using an external API.
func GetCountryFromIP(ip string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	apiURL := fmt.Sprintf("https://ipinfo.io/%s/json", ip) // Example API (replace with your provider)

	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch geo data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("geo API returned unexpected status: %d", resp.StatusCode)
	}

	var geoData GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
		return "", fmt.Errorf("failed to parse geo data: %v", err)
	}

	if geoData.Country == "" {
		return "unknown", nil
	}

	return geoData.Country, nil
}

// LogAuditEntry logs an audit entry into the audit_logs table.
func LogAuditEntry(ctx context.Context, client *dynamodb.Client, query, clientIP, country string) error {
	now := time.Now().UTC()
	entry := map[string]types.AttributeValue{
		"timestamp": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
		"query":     &types.AttributeValueMemberS{Value: query},
		"ip":        &types.AttributeValueMemberS{Value: clientIP},
		"country":   &types.AttributeValueMemberS{Value: country},
	}

	tableName := os.Getenv("AUDIT_LOGS_TABLE")
	if tableName == "" {
		tableName = "audit_logs" // Default table name
	}

	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      entry,
	})
	if err != nil {
		log.Printf("Error logging audit entry: %v", err)
		return err
	}

	log.Printf("Successfully logged audit entry for query: %s, IP: %s", query, clientIP)
	return nil
}

// GetFilteredLogs fetches logs from the last 'minutes' specified, defaulting to 24 hours if 'minutes' is 0.
func GetFilteredLogs(ctx context.Context, client *dynamodb.Client, minutes int) ([]map[string]interface{}, error) {
	var filterTime string
	if minutes > 0 {
		filterTime = time.Now().Add(-time.Duration(minutes) * time.Minute).Format(time.RFC3339)
	} else {
		// Default to the last 24 hours
		filterTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}

	tableName := os.Getenv("AUDIT_LOGS_TABLE")
	if tableName == "" {
		tableName = "audit_logs" // Default table name
	}

	log.Printf("Fetching audit logs starting from: %s in table: %s", filterTime, tableName)

	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("#ts >= :filterTime"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":filterTime": &types.AttributeValueMemberS{Value: filterTime},
		},
	}

	result, err := client.Scan(ctx, input)
	if err != nil {
		log.Printf("Error fetching filtered logs: %v", err)
		return nil, err
	}

	var logs []map[string]interface{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &logs)
	if err != nil {
		log.Printf("Error unmarshalling filtered logs: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched %d audit logs starting from %s", len(logs), filterTime)
	return logs, nil
}
