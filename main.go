package main

import (
	"log"
	"os"
	"time"

	"demodiqit_api/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Tải biến môi trường từ file .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	// 1. Khởi tạo kết nối DB
	config.ConnectDatabase()

	// 2. Khởi tạo Gin Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "pong",
			"db_status": "connected",
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		sqlDB, err := config.DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "fail", "message": "Database disconnected"})
			return
		}

		// Ping thực tế tới Neon để kiểm tra độ trễ
		start := time.Now()
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "fail", "message": "Database unreachable"})
			return
		}
		latency := time.Since(start)

		c.JSON(200, gin.H{
			"status":  "healthy",
			"latency": latency.String(),
			"db":      "Neon Postgres Connected",
		})
	})

	// Routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", controllers.Login)
	}

	// 3. Chạy server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server đang chạy tại http://localhost:%s", port)
	r.Run(":" + port)
}
