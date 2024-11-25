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
	return err
}

// GetAuditLogs fetches the audit logs from the last 24 hours
func GetAuditLogs(ctx context.Context, client *dynamodb.Client) ([]map[string]interface{}, error) {
	now := time.Now().UTC()
	lastDay := now.Add(-24 * time.Hour).Format(time.RFC3339)

	log.Printf("Fetching audit logs from %s onwards", lastDay)

	input := &dynamodb.ScanInput{
		TableName:        aws.String("audit_logs"),
		FilterExpression: aws.String("#ts >= :lastDay"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp", // Use #ts to represent the reserved keyword
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":lastDay": &types.AttributeValueMemberS{Value: lastDay},
		},
	}

	result, err := client.Scan(ctx, input)
	if err != nil {
		log.Printf("Error during DynamoDB scan: %v", err)
		return nil, err
	}

	log.Printf("Raw scan result: %+v", result)

	var logs []map[string]interface{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &logs)
	if err != nil {
		log.Printf("Error unmarshalling scan result: %v", err)
		return nil, err
	}

	log.Printf("Successfully fetched %d audit logs", len(logs))
	return logs, nil
}
