package models

import (
	"gorm.io/gorm"
)

// User represents the 'users' table in the database
type User struct {
	gorm.Model           // Contains ID, CreatedAt, UpdatedAt, DeletedAt
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	FullName string `gorm:"column:full_name"`
	Age      int    
	Active   bool   `gorm:"default:true"`
}
