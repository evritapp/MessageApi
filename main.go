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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // Default to development if not set
	}

	 err2 := godotenv.Overload(fmt.Sprintf(".env.%s", env))
	if err2 != nil {
		log.Fatalf("Error loading .env.%s file", env)
	}

	// Initialize database connection
	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connected to db")
	defer db.CloseDB()

	// Set Gin mode based on environment
	if env == "production" {
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
