package main

import (
	"fmt"
	"log"
	"time"

	"demodiqit_api/config"
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
}
