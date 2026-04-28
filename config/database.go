package config

import (
	"log"
	"demodiqit_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối đến cơ sở dữ liệu: %v\n", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Lỗi Migration: %v\n", err)
	}

	log.Println("✅ Đã kết nối và tự động tạo bảng thành công!")
}
