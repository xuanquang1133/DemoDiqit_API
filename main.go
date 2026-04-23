package main

import (
	"log"
	"os"

	"demodiqit_api/config"
	"demodiqit_api/controllers"

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

	// Khởi tạo Admin User nếu chưa có
	config.SeedAdmin()

	// 2. Khởi tạo Gin Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "pong",
			"db_status": "connected",
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
