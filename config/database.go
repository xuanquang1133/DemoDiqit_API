package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// 1. Cấu hình Logger chuyên nghiệp (Lọc query chậm > 200ms)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Query nào chậm hơn 0.2s mới log
			LogLevel:                  logger.Warn,            // Chỉ hiện Warning hoặc Error
			IgnoreRecordNotFoundError: true,                   // Bỏ qua log khi không tìm thấy record
			Colorful:                  true,
		},
	)

	// 2. Kết nối với Logger mới
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("❌ Lỗi kết nối Database: ", err)
	}

	// 3. Cấu hình Connection Pool (Tối ưu tài nguyên)
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("❌ Không thể lấy đối tượng sql.DB")
	}

	// Tránh lãng phí: Giới hạn số lượng kết nối mở đồng thời
	sqlDB.SetMaxIdleConns(5)                  // Giữ tối đa 5 kết nối rảnh
	sqlDB.SetMaxOpenConns(20)                 // Tối đa 20 kết nối mở (với CMS nội bộ vậy là quá đủ)
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Thay mới kết nối sau 30 phút để tránh lỗi 'stale connection'

	log.Println("✅ Database đã được tối ưu hiệu năng!")
	DB = database

	// Auto Migrate
	err = DB.AutoMigrate(&models.User{})
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
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db
	log.Println("Successfully connected to PostgreSQL via GORM")
}
