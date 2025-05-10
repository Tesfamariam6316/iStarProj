package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hulupay/istar-api/config"
	"github.com/hulupay/istar-api/internal/api"
	"github.com/hulupay/istar-api/internal/client"
	"github.com/hulupay/istar-api/internal/handlers"
	"github.com/hulupay/istar-api/internal/middleware"
	"github.com/hulupay/istar-api/internal/repositories"
	"github.com/hulupay/istar-api/internal/services"
	"github.com/hulupay/istar-api/pkg/logging"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           iStar API
// @version         1.0
// @description     This is the API documentation for iStar API.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        API-Key
// @description                 API Key Authentication
//
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	cfg := config.Load()
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Failed to initialize zap logger", zap.Error(err))
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	//set up gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(logging.LoggerMiddleware(sugar))
	router.Use(middleware.ErrorHandler(logger))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	istarClient := client.NewIStarClient(cfg.IStarConfigVar, logger)
	orderRepo := repositories.NewOrderRepository( /*db.Pool,*/ logger)
	orderService := services.NewOrderService(orderRepo, istarClient, logger)

	starHandler := handlers.NewStarHandler(orderService, istarClient, logger)
	premiumHandler := handlers.NewPremiumHandler(orderService, istarClient, logger)
	walletHandler := handlers.NewWalletHandler(istarClient, logger)
	webhookHandler := handlers.NewWebhookHandler(orderRepo, cfg.WebhookSecret, logger)

	router = api.SetupRouter(router, starHandler, premiumHandler, walletHandler, webhookHandler)

	// Register health check endpoint
	router.GET("/health", healthCheck)

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown setup
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("port", "8080"))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}

// HealthCheck godoc
// @Summary      Show the status of server
// @Description  Retrieve the current status of the server
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]

// FindAllResources godoc
// @Summary      Retrieve all resources
// @Description  Get a complete list of all resources managed by the server
// @Tags         resources
// @Accept       json
// @Produce      json
// @Success      200  {array}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /resources [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Add a placeholder route for finding all resources
func FindAllResources(c *gin.Context) {
	// Dummy logic, to be replaced with actual implementation.
	c.JSON(http.StatusOK, []interface{}{
		map[string]interface{}{"id": 1, "name": "Resource1"},
		map[string]interface{}{"id": 2, "name": "Resource2"},
	})
}
