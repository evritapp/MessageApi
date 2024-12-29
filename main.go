package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"messageapi.e-vrit.co.il/routes"
	"messageapi.e-vrit.co.il/utils"
	"net/http"
	"os"
)

var env string

func main() {

	utils.InitEnvVars()

	// Set Gin mode based on environment
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Define routes
	routes.SmsRoutes(router)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9092" // Default port if not specified
	}

	fmt.Printf("Server is running on port %s in %s mode\n", port, env)
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
