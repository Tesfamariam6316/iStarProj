package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	//set up gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start server
	router.Run(":8080")
}
