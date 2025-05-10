package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/hulupay/istar-api/internal/models"
	"github.com/hulupay/istar-api/internal/repositories"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WebhookHandler handles webhook events
type WebhookHandler struct {
	repo          repositories.OrderRepository
	webhookSecret string
	logger        *zap.Logger
}

// NewWebhookHandler godocs
// @Summary      Create a new webhook handler
// @Description  Initializes a new WebhookHandler
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        repo     path      repositories.OrderRepository  true  "Order repository"
// @Param        secret   path      string                       true  "Webhook secret"
// @Param        logger   path      *zap.Logger                  true  "Logger"
// @Success      200      {object}  *WebhookHandler
// @Failure      400      {object}  models.ErrorResponse
// @Router       /webhook [post]
func NewWebhookHandler(repo repositories.OrderRepository, secret string, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		repo:          repo,
		webhookSecret: secret,
		logger:        logger.Named("webhook_handler"),
	}
}

// HandleWebhookHandler godoc
// @Summary      Handle webhook events
// @Description  Handles webhook events from iStar
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        payload  body      models.WebhookPayload  true  "Webhook payload"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  models.ErrorResponse
func (h *WebhookHandler) HandleWebhookHandler(c *gin.Context) {
	if h.webhookSecret != "" {
		signature := c.GetHeader("X-iStar-Signature")
		body, err := c.GetRawData()
		if err != nil {
			h.logger.Error("Failed to read webhook body", zap.Error(err))
			c.Error(models.InternalServerError("Failed to read webhook body"))
			return
		}
		mac := hmac.New(sha256.New, []byte(h.webhookSecret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(signature), []byte(expected)) {
			h.logger.Warn("Invalid webhook signature")
			c.Error(models.UnauthorizedError("Invalid webhook signature"))
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	var payload models.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Error("Invalid webhook payload", zap.Error(err))
		c.Error(models.ValidationError("Invalid webhook payload"))
		return
	}

	orderID, ok := payload.Order["id"].(string)
	if !ok {
		h.logger.Error("Missing order ID in webhook payload")
		c.Error(models.ValidationError("Missing order ID"))
		return
	}

	status, ok := payload.Order["status"].(string)
	if !ok {
		h.logger.Error("Missing status in webhook payload")
		c.Error(models.ValidationError("Missing status"))
		return
	}

	var txHash *string
	if payload.TxHash != nil {
		th := *payload.TxHash
		txHash = &th
	}

	var completedAt *time.Time
	if payload.CompletedAt != nil {
		completedAt = payload.CompletedAt
	}

	var errorMessage *string
	if em, ok := payload.Order["error"].(string); ok {
		errorMessage = &em
	}

	err := h.repo.UpdateOrderStatus(c.Request.Context(), orderID, models.OrderStatus(status), txHash, completedAt, errorMessage)
	if err != nil {
		h.logger.Error("Failed to update order", zap.Error(err))
		c.Error(models.InternalServerError("Failed to update order"))
		return
	}

	h.logger.Info("Webhook processed",
		zap.String("event_type", payload.EventType),
		zap.String("order_id", orderID))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

/*
func VerifyWebhookSignature(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			c.Next()
			return
		}

		signature := c.GetHeader("X-iStar-Signature")
		body, _ := c.GetRawData()

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature), []byte(expected)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid webhook signature",
			})
			return
		}

		c.Set("rawBody", body)
		c.Next()
	}
}

func HandleWebhook(c *gin.Context) {
	var payload models.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Process different event types
	switch payload.EventType {
	case "order.completed":
		handleOrderCompleted(c, payload)
	case "order.failed":
		handleOrderFailed(c, payload)
	default:
		c.JSON(http.StatusOK, gin.H{"status": "unhandled_event"})
	}
}

*/
