package main

import (
	"fmt"
	"log"
	"time"

	"demodiqit_api/config"
	"demodiqit_api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database connection
	config.ConnectDB(cfg)

	// 3. Test connection with a basic query
	var currentTime time.Time
	err := config.DB.Raw("SELECT NOW();").Scan(&currentTime).Error
	if err != nil {
		log.Fatalf("Error executing query: %v\n", err)
	}

	fmt.Printf("Current time from database: %v\n", currentTime)

	// 4. Initialize Gin
	r := gin.Default()

	// 5. Apply CORS middleware
	r.Use(middleware.CorsConfig(cfg))

	// 6. Define a simple route for testing
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
			"time":   currentTime,
		})
	})

	// 7. Run the server
	fmt.Println("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
