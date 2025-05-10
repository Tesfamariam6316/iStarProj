package middleware

import (
	"github.com/hulupay/istar-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			logger.Error("Request processing error",
				zap.String("path", c.FullPath()),
				zap.Error(err))
			switch e := err.(type) {
			case *models.APIError:
				c.JSON(e.Code, gin.H{"error": e.Message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}
	}
}
