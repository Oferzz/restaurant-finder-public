package middleware

import (
	"fmt"
	"log"
	"strings"

	"server/services"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AuditLog(client *dynamodb.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/healthz" ||
			strings.HasPrefix(c.Request.URL.Path, "/admin/logs") {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		country := c.GetHeader("Country")
		query := c.Request.URL.Query()
		queryString := query.Encode()

		fmt.Printf("Captured Query: %s\n", queryString)

		// Log to the audit storage
		err := services.LogAuditEntry(c.Request.Context(), client, queryString, clientIP, country)
		if err != nil {
			log.Printf("Failed to log audit entry: %v", err)
		}

		c.Next()
	}
}
