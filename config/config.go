package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL không được cấu hình")
	}

	return &Config{
		DatabaseURL: dsn,
	}
}
