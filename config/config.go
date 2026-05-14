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
	if feOrigin == "" {
		feOrigin = "http://localhost:5173"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key" // fallback
	}

	jwtExpDaysStr := os.Getenv("JWT_EXPIRATION_DAYS")
	jwtExpDays := 7 // default to 7 days
	if jwtExpDaysStr != "" {
		if val, err := strconv.Atoi(jwtExpDaysStr); err == nil {
			jwtExpDays = val
		}
	}

	return &Config{
		DatabaseURL:       dsn,
		FrontendOrigin:    feOrigin,
		JWTSecret:         jwtSecret,
		JWTExpirationDays: jwtExpDays,
	}
}
