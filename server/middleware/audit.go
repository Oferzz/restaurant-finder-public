package middleware

import (
	"log"
	"strings"

	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AuditLog(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health checks and logs
		if c.Request.URL.Path == "/healthz" ||
			strings.HasPrefix(c.Request.URL.Path, "/admin/logs") {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		queryString := c.Request.URL.Query().Encode()

		// Fetch country for the client IP
		country, err := services.GetCountryFromIP(clientIP)
		if err != nil {
			log.Printf("Error fetching country for IP %s: %v", clientIP, err)
			country = "unknown"
		}

		// Log to audit storage
		err = services.LogAuditEntry(c.Request.Context(), client, queryString, clientIP, country)
		if err != nil {
			log.Printf("Failed to log audit entry: %v", err)
		}

		c.Next()
	}
}
