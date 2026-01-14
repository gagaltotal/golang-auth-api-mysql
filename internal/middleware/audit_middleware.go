package middleware

import (
	"bytes"
	"encoding/json"
	"golang-auth-api-mysql/internal/domain"
	"golang-auth-api-mysql/internal/service"
	"golang-auth-api-mysql/pkg/jwt"
	"io"

	"github.com/gin-gonic/gin"
)

func AuditMiddleware(auditService service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip audit for public endpoints
		if c.Request.Method == "GET" && c.Request.URL.Path == "/api/v1/products" {
			c.Next()
			return
		}

		// Get user from context
		claims, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		userClaims := claims.(*jwt.Claims)

		// Read request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create audit log
		audit := &domain.AuditLog{
			UserID:    userClaims.UserID,
			Action:    c.Request.Method,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		// Add payload summary (max 500 chars)
		if len(bodyBytes) > 0 {
			var payload map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &payload); err == nil {
				// Remove sensitive fields
				delete(payload, "password")
				delete(payload, "old_password")
				delete(payload, "new_password")

				payloadJSON, _ := json.Marshal(payload)
				if len(payloadJSON) > 500 {
					audit.Payload = string(payloadJSON[:500]) + "..."
				} else {
					audit.Payload = string(payloadJSON)
				}
			}
		}

		// Determine resource
		if len(c.Params) > 0 {
			audit.Resource = c.Param("id")
		}

		// Save audit log asynchronously
		go auditService.Create(audit)

		c.Next()
	}
}
