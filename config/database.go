package config

import (
	"log"
	"os"

	"demodiqit_api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db
	log.Println("Successfully connected to PostgreSQL via GORM")

	// Auto Migrate
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migration completed")
}

func SeedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("email = ?", "admin@gmail.com").Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password for admin seed")
		}

		admin := models.User{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		result := DB.Create(&admin)
		if result.Error != nil {
			log.Fatal("Failed to seed admin user: ", result.Error)
		}
		log.Println("Admin user seeded successfully")
	} else {
		log.Println("Admin user already exists")
	}
}
