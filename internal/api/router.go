package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hulupay/istar-api/internal/handlers"
)

func SetupRouter(
	route *gin.Engine,
	starHandler *handlers.StarHandler,
	premiumHandler *handlers.PremiumHandler,
	walletHandler *handlers.WalletHandler,
	webhookHandler *handlers.WebhookHandler) *gin.Engine {

	// Star Gifting
	route.GET("/star/recipient/search", starHandler.SearchStarRecipientHandler)
	route.POST("/orders/star", starHandler.CreateStarGiftAsyncHandler)
	route.POST("/orders/star/sync", starHandler.CreateStarGiftSyncHandler)

	// Premium Gifts
	route.GET("/premium/recipient/search", premiumHandler.SearchPremiumRecipientHandler)
	route.POST("/orders/premium", premiumHandler.CreatePremiumGiftAsyncHandler)
	route.POST("/orders/premium/sync", premiumHandler.CreatePremiumGiftSyncHandler)
	route.GET("/premium/packages", premiumHandler.GetPremiumPackagesHandler)

	// Wallet
	route.GET("/wallet/balance", walletHandler.GetWalletBalanceHandler)

	// Webhooks
	route.POST("/webhooks/istar", webhookHandler.HandleWebhookHandler)

	return route
}
