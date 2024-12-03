package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"messageapi.e-vrit.co.il/routes"
)

var env string

func main() {
	fmt.Printf("runnin in version %v", env)

	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		log.Fatalf("Error loading .env.%s file", env)
	}
	// env := os.Getenv("GO_ENV")

	// Initialize database connection
	// err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("connected to db")
	// defer db.CloseDB()

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
		port = "8080" // Default port if not specified
	}
	if os.Getenv("ASPNETCORE_PORT") != "" { // get enviroment variable that set by ACNM
		port = os.Getenv("ASPNETCORE_PORT")
	}
	fmt.Printf("Server is running on port %s in %s mode\n", port, env)
	router.Run(":" + port)
}
