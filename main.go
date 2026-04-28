package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Tải các biến môi trường từ file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL không được cấu hình")
	}

	// 2. Khởi tạo kết nối với CSDL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối đến cơ sở dữ liệu: %v\n", err)
	}

	log.Println("✅ Kết nối đến Neon Postgres thành công!")

	// 3. Kiểm tra (Testing) - Thực hiện truy vấn đơn giản
	var currentTime time.Time
	err = db.Raw("SELECT NOW();").Scan(&currentTime).Error
	if err != nil {
		log.Fatalf("Lỗi khi thực hiện truy vấn: %v\n", err)
	}

	fmt.Printf("Thời gian hiện tại từ cơ sở dữ liệu: %v\n", currentTime)
}
