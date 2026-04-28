package models

import (
	"gorm.io/gorm"
)

// User đại diện cho bảng 'users' trong database
type User struct {
	gorm.Model           // Chứa ID, CreatedAt, UpdatedAt, DeletedAt
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	FullName string `gorm:"column:full_name"`
	Age      int    
	Active   bool   `gorm:"default:true"`
}
