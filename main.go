package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"demodiqit_api/config"
	"demodiqit_api/middleware"
	"demodiqit_api/routes"

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

	// Setup API routes
	api := r.Group("/api/v1")
	routes.SetupAuthRoutes(api, cfg)

	api.Use(middleware.JWTAuthMiddleware(cfg))
	routes.SetupUserRoutes(api, cfg)

	// 7. Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running at http://localhost:%s\n", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
