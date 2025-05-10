package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func RequireHTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS == nil && c.Request.URL.Scheme != "https" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "HTTPS required",
			})
			return
		}
		c.Next()
	}
}

func VerifyWebhookSignature(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-iStar-Signature")
		body, _ := c.GetRawData()

		// Verify HMAC-SHA256 signature
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expected)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid webhook signature",
			})
			return
		}

		// Restore the body for subsequent handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
