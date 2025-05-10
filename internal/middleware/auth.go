package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func APIKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*zap.Logger)

		apiKey := GetAPIKey(c)
		if apiKey == "" {
			logger.Warn("Missing API key")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
				"code":  "MISSING_API_KEY",
			})
			return
		}

		if !isValidAPIKey(apiKey, validKey) {
			logger.Warn("Invalid API key attempt", zap.String("key", apiKey))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
				"code":  "INVALID_API_KEY",
			})
			return
		}

		c.Next()
	}
}

// GetAPIKey extracts and sanitizes the API key from headers
func GetAPIKey(c *gin.Context) string {
	return strings.TrimSpace(c.GetHeader("API-Key"))
}

// isValidAPIKey securely compares keys using constant time comparison
func isValidAPIKey(inputKey, validKey string) bool {
	if inputKey == "" || validKey == "" {
		return false
	}
	return subtleConstantTimeCompare(inputKey, validKey)
}

// subtleConstantTimeCompare prevents timing attacks
func subtleConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
