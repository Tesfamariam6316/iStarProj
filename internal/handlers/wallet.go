package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/hulupay/istar-api/internal/client"
	"go.uber.org/zap"
	"net/http"
)

// WalletHandler handles wallet-related endpoints
type WalletHandler struct {
	istarClient *client.IStarClient
	logger      *zap.Logger
}

// NewWalletHandler initializes a new WalletHandler
// GetWalletBalanceHandler godoc
// @Summary      Retrieve wallet balance
// @Description  Retrieves the wallet balance of the current user
// @Tags         wallet
// @Produce      json
// @Success      200    {object}  map[string]interface{}
// @Router       /wallet/balance [get]
func NewWalletHandler(istarClient *client.IStarClient, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{
		istarClient: istarClient,
		logger:      logger.Named("wallet_handler"),
	}
}

// GetWalletBalanceHandler godoc
// @Summary      Retrieve wallet balance
// @Description  Retrieves the wallet balance of the current user
// @Tags         wallet
// @Produce      json
// @Success      200    {object}  map[string]interface{}
// @Router       /wallet/balance [get]
func (h *WalletHandler) GetWalletBalanceHandler(c *gin.Context) {
	ctx := c.Request.Context()
	resp, err := h.istarClient.DoRequest(ctx, "GET", "/wallet/balance", nil)
	if err != nil {
		h.logger.Error("Failed to retrieve wallet balance", zap.Error(err))
		c.Error(err)
		return
	}

	h.logger.Info("Wallet balance retrieved")
	c.JSON(http.StatusOK, resp)
}
