package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"messageapi.e-vrit.co.il/db"
	"messageapi.e-vrit.co.il/routes"
)

func main() {
	err := godotenv.Load()
	env := os.Getenv("GO_ENV")
	// Initialize database connection
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connected to db")
	defer db.CloseDB()

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
		port = "8080" // Default port if not specified
	}
	fmt.Printf("Server is running on port %s in %s mode\n", port, env)
	router.Run(":" + port)

}
