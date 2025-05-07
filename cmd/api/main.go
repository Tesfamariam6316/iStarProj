package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hulupay/istar-api/pkg/logging"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
)

// @title           iStar API
// @version         1.0
// @description     This is the API documentation for iStar API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        API-Key
// @description                 API Key Authentication

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Register health check endpoint
	router.GET("/health", healthCheck)

	// Start server
	router.Run(":8080")
}

// HealthCheck godoc
// @Summary      Show the status of server
// @Description  get the status of server
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
