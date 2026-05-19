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

	err = DB.AutoMigrate(&models.User{}, &models.Product{})
	if err != nil {
		log.Fatalf("Migration Error: %v\n", err)
	}

	log.Println("✅ Connected and tables auto-migrated successfully!")

	// Seed Admin User
	var admin models.User
	if err := DB.Unscoped().Where("username = ?", "admin").First(&admin).Error; err != nil {
		// Create new admin if it doesn't exist
		newAdmin := models.User{
			Username: "admin",
			Password: "123456", // BeforeCreate hook will hash this!
			Email:    "admin@gmail.com",
			FullName: "Administrator",
		}
		if err := DB.Create(&newAdmin).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}
		log.Println("✅ Admin user (admin@gmail.com / 123456) seeded successfully!")
	} else {
		// Restore if soft deleted
		if admin.DeletedAt.Valid {
			DB.Unscoped().Model(&admin).Update("deleted_at", nil)
			log.Println("✅ Admin user was soft-deleted, restored successfully!")
		}
		// Update email if the user already exists (e.g. from a previous seed)
		if admin.Email != "admin@gmail.com" {
			admin.Email = "admin@gmail.com"
			DB.Save(&admin)
			log.Println("✅ Admin user email updated to admin@gmail.com!")
		}
	}
}
