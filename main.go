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
	// 1. Load cấu hình
	cfg := config.LoadConfig()

	// 2. Khởi tạo kết nối Database
	config.ConnectDB(cfg)

	// 3. Kiểm tra kết nối bằng truy vấn cơ bản
	var currentTime time.Time
	err := config.DB.Raw("SELECT NOW();").Scan(&currentTime).Error
	if err != nil {
		log.Fatalf("Lỗi khi thực hiện truy vấn: %v\n", err)
	}

	fmt.Printf("Thời gian hiện tại từ cơ sở dữ liệu: %v\n", currentTime)

	// 4. Khởi tạo Gin
	r := gin.Default()

	// 5. Áp dụng CORS middleware
	r.Use(middleware.CorsConfig(cfg))

	// 6. Định nghĩa route đơn giản để test
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "UP",
			"time":   currentTime,
		})
	})

	// 7. Chạy server
	fmt.Println("Server đang chạy tại http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
