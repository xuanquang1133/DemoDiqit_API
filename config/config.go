package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	FrontendOrigin string
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

	return &Config{
		DatabaseURL:    dsn,
		FrontendOrigin: feOrigin,
	}
}
