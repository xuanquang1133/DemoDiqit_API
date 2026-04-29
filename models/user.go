package models

import (
	"gorm.io/gorm"
)

// User đại diện cho bảng 'users' trong database
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Không trả về password
	Email    string `gorm:"unique;not null" json:"email"`
	FullName string `json:"full_name"`
	Role     string `gorm:"default:customer" json:"role"` // Thêm phân quyền
}
