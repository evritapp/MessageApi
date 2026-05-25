package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"messageapi.e-vrit.co.il/routes"
	"messageapi.e-vrit.co.il/utils"
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
