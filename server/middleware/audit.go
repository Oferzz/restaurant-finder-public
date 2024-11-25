package middleware

import (
	"log"

	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AuditLog(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture request details
		clientIP := c.ClientIP()
		country := c.GetHeader("X-Country") // Assuming the client sends this header
		query := c.Request.URL.Query()

		// Log to the audit storage
		err := services.LogAuditEntry(c.Request.Context(), client, query.Encode(), clientIP, country)
		if err != nil {
			log.Printf("Failed to log audit entry: %v", err)
		}

		c.Next()
	}
}
