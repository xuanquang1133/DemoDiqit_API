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
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Migration Error: %v\n", err)
	}

	log.Println("✅ Connected and tables auto-migrated successfully!")
}
