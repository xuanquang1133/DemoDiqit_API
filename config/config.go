package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	FrontendOrigin    string
	JWTSecret         string
	JWTExpirationDays int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	feOrigin := os.Getenv("FE_ORIGIN")

	jwtSecret := os.Getenv("JWT_SECRET")

	jwtExpDaysStr := os.Getenv("JWT_EXPIRATION_DAYS")
	jwtExpDays, _ := strconv.Atoi(jwtExpDaysStr)

	return &Config{
		DatabaseURL:       dsn,
		FrontendOrigin:    feOrigin,
		JWTSecret:         jwtSecret,
		JWTExpirationDays: jwtExpDays,
	}
}
