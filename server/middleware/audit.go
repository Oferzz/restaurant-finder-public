package middleware

import (
	"log"
	"strings"

	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// AuditLog middleware captures request details and logs them to the audit storage.
func AuditLog(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health checks and admin log endpoints
		if c.Request.URL.Path == "/healthz" ||
			strings.HasPrefix(c.Request.URL.Path, "/admin/logs") {
			c.Next()
			return
		}

		// Extract client IP and country
		clientIP := c.ClientIP()
		country := c.GetHeader("Country")
		if country == "" {
			// Attempt to fetch the country dynamically if not provided in the headers
			var err error
			country, err = services.GetCountryFromIP(clientIP)
			if err != nil {
				log.Printf("Failed to fetch country for IP %s: %v", clientIP, err)
				country = "unknown"
			}
		}

		// Extract query parameters
		queryString := c.Request.URL.RawQuery

		log.Printf("Captured Query: %s, Client IP: %s, Country: %s", queryString, clientIP, country)

		// Log the audit entry
		err := services.LogAuditEntry(c.Request.Context(), client, queryString, clientIP, country)
		if err != nil {
			log.Printf("Failed to log audit entry: %v", err)
		}

		// Proceed to the next middleware or handler
		c.Next()
	}
}
