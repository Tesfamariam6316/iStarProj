package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hulupay/istar-api/internal/client"
	"github.com/hulupay/istar-api/internal/models"
	"github.com/hulupay/istar-api/internal/services"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// PremiumHandler handles premium gift and package endpoints
type PremiumHandler struct {
	orderService services.OrderService
	istarClient  *client.IStarClient
	logger       *zap.Logger
}

// NewPremiumHandler initializes a new PremiumHandler
// NewPremiumHandler godoc
// @Summary      Premium handler setup
// @Description  Handle operations related to premium gifting
// @Tags         premium
// @Router       /premium/recipient/search [get]
func NewPremiumHandler(orderService services.OrderService, istarClient *client.IStarClient, logger *zap.Logger) *PremiumHandler {
	return &PremiumHandler{
		orderService: orderService,
		istarClient:  istarClient,
		logger:       logger.Named("premium_handler"),
	}
}

// SearchPremiumRecipientHandler godoc
// @Summary      Search for a premium recipient
// @Description  Searches for a premium recipient by username and months
// @Tags         premium
// @Accept       json
// @Produce      json
// @Param        username  query     string  true  "Username of the recipient"
// @Param        months    query     int     true  "Number of months (3, 6, or 12)"
// @Success      200       {object}  models.PremiumRecipientResponse
// @Failure      400       {object}  models.ErrorResponse
func (h *PremiumHandler) SearchPremiumRecipientHandler(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Query("username")
	monthsStr := c.Query("months")

	if username == "" || monthsStr == "" {
		h.logger.Error("Missing required parameters")
		c.Error(models.ValidationError("Missing username or months"))
		return
	}

	months, err := strconv.Atoi(monthsStr)
	if err != nil || !isValidMonths(months) {
		h.logger.Error("Invalid months")
		c.Error(models.ValidationError("Months must be 3, 6, or 12"))
		return
	}

	resp, err := h.istarClient.DoRequest(ctx, "GET", fmt.Sprintf("/premium/recipient/search?username=%s&months=%d", username, months), nil)
	if err != nil {
		h.logger.Error("Failed to search premium recipient", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Premium recipient searched", zap.String("username", username))
	c.JSON(http.StatusOK, resp)
}

// CreatePremiumGiftAsyncHandler godoc
// @Summary      Create a premium gift order (asynchronous)
// @Description  Creates a premium gift order asynchronously
// @Tags         premium
// @Accept       json
// @Produce      json
// @Param        request  body     models.CreatePremiumOrderRequest  true  "Create premium order request"
// @Success      202      {object}  models.CreatePremiumOrderResponse
// @Failure      400      {object}  models.ErrorResponse
func (h *PremiumHandler) CreatePremiumGiftAsyncHandler(c *gin.Context) {
	var req models.CreatePremiumOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.Error(models.ValidationError("Invalid request body: " + err.Error()))
		return
	}

	if req.Username == "" || req.RecipientHash == "" || !isValidMonths(req.Months) || req.WalletType == "" {
		h.logger.Error("Invalid request parameters")
		c.Error(models.ValidationError("Invalid request parameters: username, recipient_hash, months (3, 6, 12), wallet_type required"))
		return
	}

	resp, err := h.orderService.CreatePremiumOrderAsync(c, req)
	if err != nil {
		h.logger.Error("Failed to create premium gift order", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Premium gift order created (async)", zap.String("order_id", resp.ID.String()))
	c.JSON(http.StatusAccepted, resp)
}

// CreatePremiumGiftSyncHandler godoc
// @Summary      Create a premium gift order (synchronous)
// @Description  Creates a premium gift order synchronously
// @Tags         premium
// @Accept       json
// @Produce      json
// @Param        request  body     models.CreatePremiumOrderRequest  true  "Create premium order request"
// @Success      200      {object}  models.CreatePremiumOrderResponse
// @Failure      400      {object}  models.ErrorResponse
func (h *PremiumHandler) CreatePremiumGiftSyncHandler(c *gin.Context) {
	var req models.CreatePremiumOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.Error(models.ValidationError("Invalid request body: " + err.Error()))
		return
	}

	if req.Username == "" || req.RecipientHash == "" || !isValidMonths(req.Months) || req.WalletType == "" {
		h.logger.Error("Invalid request parameters")
		c.Error(models.ValidationError("Invalid request parameters: username, recipient_hash, months (3, 6, 12), wallet_type required"))
		return
	}

	resp, err := h.orderService.CreatePremiumOrderSync(c, req)
	if err != nil {
		h.logger.Error("Failed to create premium gift order", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Premium gift order created (sync)", zap.String("order_id", resp.ID.String()))
	c.JSON(http.StatusOK, resp)
}

// GetPremiumPackagesHandler godoc
// @Summary      Retrieve premium packages
// @Description  Retrieves the available premium packages
// @Tags         premium
// @Produce      json
// @Success      200      {object}  models.PremiumPackagesResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /premium/packages [get]
func (h *PremiumHandler) GetPremiumPackagesHandler(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := h.istarClient.DoRequest(ctx, "GET", "/premium/packages", nil)
	if err != nil {
		h.logger.Error("Failed to retrieve premium packages", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Premium packages retrieved")
	c.JSON(http.StatusOK, resp)
}

// isValidMonths checks if the given months value is valid (3, 6, or 12)
func isValidMonths(months int) bool {
	return months == 3 || months == 6 || months == 12
}

/*

// SearchPremiumRecipient Search Premium Recipient
func SearchPremiumRecipient(logger *zap.Logger, baseURL, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.Query("username")
		months := c.Query("months")

		if username == "" || months == "" {
			logger.Error("Missing required parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or months"})
			return
		}

		url := baseURL + "/premium/recipient/search?username=" + username + "&months=" + months
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			logger.Error("Failed to create request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("API-Key", apiKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error("Failed to send request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logger.Error("Failed to parse response", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		if resp.StatusCode >= 400 {
			logger.Error("iStar API error", zap.Any("error", result["error"]))
			c.JSON(resp.StatusCode, result)
			return
		}

		logger.Info("Premium recipient searched", zap.String("username", username))
		c.JSON(resp.StatusCode, result)
	}
}

// CreatePremiumOrder Create Premium Order (Asynchronous):
func CreatePremiumOrder(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username      string `json:"username"`
		RecipientHash string `json:"recipient_hash"`
		Months        int    `json:"months"`
		WalletType    string `json:"wallet_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		main.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.RecipientHash == "" || req.Months <= 0 || req.WalletType == "" {
		main.logger.Error("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	url := main.baseURL + "/orders/premium"
	payload, err := json.Marshal(req)
	if err != nil {
		main.logger.Error("Failed to marshal request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		main.logger.Error("Failed to create request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	httpReq.Header.Set("API-Key", main.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		main.logger.Error("Failed to send request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		main.logger.Error("Failed to parse response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	if resp.StatusCode >= 400 {
		main.logger.Error("iStar API error", zap.Any("error", result["error"]))
		c.JSON(resp.StatusCode, result)
		return
	}

	orderID := uuid.New().String()
	istarOrderID := result["order_id"].(string)
	status := result["status"].(string)
	amount := result["amount"].(float64)
	createdAt := time.Now()
	months := req.Months

	_, err = db.Pool.Exec(ctx, `
        INSERT INTO orders (id, type, status, username, recipient_hash, months, amount, wallet_type, created_at, updated_at)
        VALUES ($1, 'premium', $2, $3, $4, $5, $6, $7, $8, $8)`,
		orderID, status, req.Username, req.RecipientHash, months, amount, req.WalletType, createdAt)
	if err != nil {
		main.logger.Error("Failed to insert order", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert order"})
		return
	}

	main.logger.Info("Premium order created", zap.String("order_id", orderID))
	c.JSON(resp.StatusCode, result)
}

// CreatePremiumOrderSync Create Premium Order (Synchronous):
func CreatePremiumOrderSync(c *gin.Context) {
	ctx := c.Request.Context()
	var req struct {
		Username      string `json:"username"`
		RecipientHash string `json:"recipient_hash"`
		Months        int    `json:"months"`
		WalletType    string `json:"wallet_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		main.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" || req.RecipientHash == "" || req.Months <= 0 || req.WalletType == "" {
		main.logger.Error("Missing required fields")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	url := main.baseURL + "/orders/premium/sync"
	payload, err := json.Marshal(req)
	if err != nil {
		main.logger.Error("Failed to marshal request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		main.logger.Error("Failed to create request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	httpReq.Header.Set("API-Key", main.apiKey

func GetPremiumPackages(context *gin.Context) {

}

	func
	GetPremiumPackages(context * gin.Context)
	{

	}


*/
