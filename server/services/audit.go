package services

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// LogAuditEntry logs an audit entry into the audit_logs table
func LogAuditEntry(ctx context.Context, client *dynamodb.Client, query, clientIP, country string) error {
	now := time.Now().UTC()
	entry := map[string]types.AttributeValue{
		"timestamp": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
		"query":     &types.AttributeValueMemberS{Value: query},
		"ip":        &types.AttributeValueMemberS{Value: clientIP},
		"country":   &types.AttributeValueMemberS{Value: country},
	}

	_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("audit_logs"),
		Item:      entry,
	})
	if err != nil {
		log.Printf("Error logging audit entry: %v", err)
		return err
	}

	log.Printf("Successfully logged audit entry for query: %s, IP: %s", query, clientIP)
	return nil
}

// GetAuditLogs fetches the audit logs from the last 24 hours
func GetAuditLogs(ctx context.Context, client *dynamodb.Client) ([]map[string]interface{}, error) {
	return GetFilteredLogs(ctx, client, 1440) // 1440 minutes = 24 hours
}

// GetFilteredLogs fetches logs from the last 'minutes' specified, defaulting to 24 hours if 'minutes' is 0
func GetFilteredLogs(ctx context.Context, client *dynamodb.Client, minutes int) ([]map[string]interface{}, error) {
	var filterTime string
	if minutes > 0 {
		filterTime = time.Now().Add(-time.Duration(minutes) * time.Minute).Format(time.RFC3339)
	} else {
		// Default to the last 24 hours
		filterTime = time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	}

	log.Printf("Fetching audit logs starting from: %s", filterTime)

	input := &dynamodb.ScanInput{
		TableName:        aws.String("audit_logs"),
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
