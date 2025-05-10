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

// StarHandler handles star gifting endpoints
type StarHandler struct {
	orderService services.OrderService
	istarClient  *client.IStarClient
	logger       *zap.Logger
}

// NewStarHandler godoc
// @Summary      Create a new star handler
// @Description  Initializes a new star handler with the given order service and iStar client
// @Tags         star
// @Accept       json
// @Produce      json
// @Param        orderService  query     services.OrderService  true  "Order service"
// @Param        istarClient   query     *client.IStarClient     true  "iStar client"
// @Param        logger        query     *zap.Logger            true  "Logger"
// @Success      200          {object}  handlers.StarHandler
// @Failure      400          {object}  models.ErrorResponse
// @Router       /star/handler [get]
// NewStarHandler initializes a new StarHandler
func NewStarHandler(orderService services.OrderService, istarClient *client.IStarClient, logger *zap.Logger) *StarHandler {
	return &StarHandler{
		orderService: orderService,
		istarClient:  istarClient,
		logger:       logger.Named("star_handler"),
	}
}

// SearchStarRecipientHandler godoc
// @Summary      Search for star recipients
// @Description  Retrieves a list of potential recipients for star gifting
// @Tags         star
// @Accept       json
// @Produce      json
// @Param        username  query     string  true  "Username to search for"
// @Param        quantity  query     int     true  "Quantity of stars to gift (50-1,000,000)"
// @Success      200       {array}   map[string]interface{}
// @Failure      400       {object}
// @Router       /star/recipient/search [get]
func (h *StarHandler) SearchStarRecipientHandler(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Query("username")
	quantityStr := c.Query("quantity")

	if username == "" || quantityStr == "" {
		h.logger.Error("Missing required parameters")
		c.Error(models.ValidationError("Missing username or quantity"))
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 50 || quantity > 1000000 {
		h.logger.Error("Invalid quantity")
		c.Error(models.ValidationError("Quantity must be between 50 and 1,000,000"))
		return
	}

	resp, err := h.istarClient.DoRequest(ctx, "GET", fmt.Sprintf("/star/recipient/search?username=%s&quantity=%d", username, quantity), nil)
	if err != nil {
		h.logger.Error("Failed to search star recipient", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Star recipient searched", zap.String("username", username))
	c.JSON(http.StatusOK, resp)
}

// CreateStarGiftAsyncHandler godoc
// @Summary      Create star gift order (asynchronous)
// @Description  Creates a star gift order asynchronously
// @Tags         star
// @Accept       json
// @Produce      json
// @Param        request  body     models.CreateStarOrderRequest  true  "Create star order request"
// @Success      202      {object}  models.CreateStarOrderResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /star/gift/async [post]
func (h *StarHandler) CreateStarGiftAsyncHandler(c *gin.Context) {
	var req models.CreateStarOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.Error(models.ValidationError("Invalid request body: " + err.Error()))
		return
	}

	if req.Username == "" || req.RecipientHash == "" || req.Quantity < 50 || req.Quantity > 1000000 || req.WalletType == "" {
		h.logger.Error("Invalid request parameters")
		c.Error(models.ValidationError("Invalid request parameters: username, recipient_hash, quantity (50-1,000,000), wallet_type required"))
		return
	}

	resp, err := h.orderService.CreateStarOrderAsync(c, req)
	if err != nil {
		h.logger.Error("Failed to create star gift order", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Star gift order created (async)", zap.String("order_id", resp.ID.String()))
	c.JSON(http.StatusAccepted, resp)
}

// CreateStarGiftSyncHandler godoc
// @Summary      Create star gift order (synchronous)
// @Description  Creates a star gift order synchronously
// @Tags         star
// @Accept       json
// @Produce      json
// @Param        request  body     models.CreateStarOrderRequest  true  "Create star order request"
// @Success      200      {object}  models.CreateStarOrderResponse
// @Failure      400      {object}  models.ErrorResponse
// @Router       /star/gift/sync [post]
func (h *StarHandler) CreateStarGiftSyncHandler(c *gin.Context) {
	var req models.CreateStarOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.Error(models.ValidationError("Invalid request body: " + err.Error()))
		return
	}

	if req.Username == "" || req.RecipientHash == "" || req.Quantity < 50 || req.Quantity > 1000000 || req.WalletType == "" {
		h.logger.Error("Invalid request parameters")
		c.Error(models.ValidationError("Invalid request parameters: username, recipient_hash, quantity (50-1,000,000), wallet_type required"))
		return
	}

	resp, err := h.orderService.CreateStarOrderSync(c, req)
	if err != nil {
		h.logger.Error("Failed to create star gift order", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Star gift order created (sync)", zap.String("order_id", resp.ID.String()))
	c.JSON(http.StatusOK, resp)
}

/*
// SearchStarRecipient godoc
// @Summary      Search for star recipients
// @Description  Retrieves a list of potential recipients for star gifting
// @Tags         star
// @Accept       json
// @Produce      json
// @Param        query  query     string  true  "Search query"
// @Success      200    {array}   map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Router       /star/recipient/search [get]
func SearchStarRecipient(logger *zap.Logger, baseURL, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.Query("username")
		quantity := c.Query("quantity")

		if username == "" || quantity == "" {
			logger.Error("Missing required parameters")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing username or quantity"})
			return
		}

		url := baseURL + "/star/recipient/search?username=" + username + "&quantity=" + quantity
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

		logger.Info("Star recipient searched", zap.String("username", username))
		c.JSON(resp.StatusCode, result)
	}
}

// CreateStarOrder Create Star Order (Asynchronous):
func CreateStarOrder(logger *zap.Logger, baseURL, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req struct {
			Username      string `json:"username"`
			RecipientHash string `json:"recipient_hash"`
			Quantity      int    `json:"quantity"`
			WalletType    string `json:"wallet_type"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if req.Username == "" || req.RecipientHash == "" || req.Quantity <= 0 || req.WalletType == "" {
			logger.Error("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		url := baseURL + "/orders/star"
		payload, err := json.Marshal(req)
		if err != nil {
			logger.Error("Failed to marshal request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			logger.Error("Failed to create request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		httpReq.Header.Set("API-Key", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(httpReq)
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

		//save to db
		orderID := uuid.New().String()

		logger.Info("Star order created", zap.String("order_id", orderID))
		c.JSON(resp.StatusCode, result)
	}
}

// CreateStarOrderSync Create Star Order (Synchronous):
func CreateStarOrderSync(logger *zap.Logger, baseURL, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req struct {
			Username      string `json:"username"`
			RecipientHash string `json:"recipient_hash"`
			Quantity      int    `json:"quantity"`
			WalletType    string `json:"wallet_type"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.Error("Invalid request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if req.Username == "" || req.RecipientHash == "" || req.Quantity <= 0 || req.WalletType == "" {
			logger.Error("Missing required fields")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		url := baseURL + "/orders/star/sync"
		payload, err := json.Marshal(req)
		if err != nil {
			logger.Error("Failed to marshal request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			logger.Error("Failed to create request", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		httpReq.Header.Set("API-Key", apiKey)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(httpReq)
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

		//save to db
		orderID := uuid.New().String()

		logger.Info("Synchronous star order created", zap.String("order_id", orderID))
		c.JSON(resp.StatusCode, result)
	}
}

*/
